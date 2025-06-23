package authcontrollers

import (
	"context"
	"curd-web-go/config"
	"curd-web-go/entities"
	"curd-web-go/helpers"
	"database/sql"
	"encoding/json"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

// POST /register
func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	tr := otel.Tracer("authcontrollers")
	ctx, span := tr.Start(ctx, "Register")
	defer span.End()

	ip := strings.Split(r.RemoteAddr, ":")[0]
	if helpers.IsRateLimited(ip) {
		helpers.ResponseJSON(w, http.StatusTooManyRequests, map[string]string{"message": "Rate limit exceeded"})
		return
	}

	var UserInput entities.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&UserInput); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	defer r.Body.Close()

	if !helpers.IsValidEmail(UserInput.Email) {
		helpers.ResponseJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid email format"})
		return
	}
	if !helpers.IsStrongPassword(UserInput.Password) {
		helpers.ResponseJSON(w, http.StatusBadRequest, map[string]string{
			"message": "Password too weak. Minimum 8 characters, with at least one number and one symbol.",
		})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(UserInput.Password), bcrypt.DefaultCost)
	if err != nil {
		span.RecordError(err)
		helpers.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to hash password"})
		return
	}
	UserInput.Password = string(hashPassword)

	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3)`
	_, err = config.DB.Exec(query, UserInput.Name, UserInput.Email, UserInput.Password)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed inserting user")
		response := map[string]string{"message": err.Error()}
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "Succes"}
	helpers.ResponseJSON(w, http.StatusOK, response)

}

// POST /login
func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	tr := otel.Tracer("authcontrollers")
	ctx, span := tr.Start(ctx, "Login")
	defer span.End()

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if helpers.IsRateLimited(ip) {
		helpers.ResponseJSON(w, http.StatusTooManyRequests, map[string]string{"message": "Rate limit exceeded"})
		return
	}

	// GET method: check token
	if r.Method == http.MethodGet {
		cookie, err := r.Cookie("Token")
		if err != nil {
			span.RecordError(err)
			helpers.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token not found"})
			return
		}

		// Validate token
		tokenStr := cookie.Value
		claims := &config.ClaimsJWT{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return config.Jwt_Secret, nil
		})
		if err != nil || !token.Valid {
			span.RecordError(err)
			helpers.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
			return
		}

		helpers.ResponseJSON(w, http.StatusOK, map[string]string{"message": "Token valid", "email": claims.Username})
		return
	}

	// POST method: login user
	if r.Method == http.MethodPost {
		var UserInput entities.User
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&UserInput); err != nil {
			span.RecordError(err)
			response := map[string]string{"message": err.Error()}
			helpers.ResponseJSON(w, http.StatusBadRequest, response)
			return
		}
		defer r.Body.Close()

		span.SetAttributes(attribute.String("user.email", UserInput.Email))

		var User entities.User
		query := `SELECT id, name, email, password FROM users WHERE email =$1`
		err := config.DB.QueryRow(query, UserInput.Email).Scan(&User.ID, &User.Name, &User.Email, &User.Password)
		if err != nil {
			span.RecordError(err)
			switch err {
			case sql.ErrNoRows:
				helpers.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"message": "Wrong Username Or Password"})
				return
			default:
				helpers.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
				return
			}
		}

		if err := bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(UserInput.Password)); err != nil {
			span.RecordError(err)
			helpers.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"message": "Wrong Username Or Password"})
			return
		}

		expiredTime := time.Now().Add(time.Minute * 5)
		refreshTime := time.Now().Add(24 * time.Hour)

		claims := &config.ClaimsJWT{
			Username: User.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "crud-web-go-Token",
				ExpiresAt: jwt.NewNumericDate(expiredTime),
			},
		}

		JwtAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		token, err := JwtAlgo.SignedString(config.Jwt_Secret)
		if err != nil {
			span.RecordError(err)
			helpers.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}

		refreshClaims := &config.ClaimsRefreshJWT{
			Username: User.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "crud-web-go-RefreshToken",
				ExpiresAt: jwt.NewNumericDate(refreshTime),
			},
		}
		refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(config.Jwt_Secret)
		if err != nil {
			span.RecordError(err)
			helpers.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to generate refresh token"})
			return
		}

		ctx := context.Background()
		if config.RedisClient == nil {
			span.RecordError(err)
			helpers.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": "Redis not initialized"})
			return
		}

		err = config.RedisClient.Set(ctx, "refresh:"+User.Email, refreshToken, time.Minute*5).Err()
		if err != nil {
			span.RecordError(err)
			helpers.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to store refresh token"})
			return
		}

		err = config.RedisClient.Set(ctx, User.Email, token, time.Minute*5).Err()
		if err != nil {
			span.RecordError(err)
			helpers.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to save token to Redis"})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "Token",
			Path:     "/",
			Value:    token,
			HttpOnly: true,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "RefreshToken",
			Path:     "/",
			Value:    token,
			HttpOnly: true,
		})

		helpers.ResponseJSON(w, http.StatusOK, map[string]string{"message": "Login Success"})
		return
	}

	helpers.ResponseJSON(w, http.StatusMethodNotAllowed, map[string]string{"message": "Method not allowed"})
}

func Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	tr := otel.Tracer("authcontrollers")
	ctx, span := tr.Start(ctx, "Logout")
	defer span.End()

	cookie, err := r.Cookie("Token")
	if err == nil && cookie != nil {
		claims := &config.ClaimsJWT{}
		tokenStr := cookie.Value
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return config.Jwt_Secret, nil
		})
		if err == nil && token.Valid {
			_ = config.RedisClient.Del(context.Background(), "access:"+claims.Username).Err()
			_ = config.RedisClient.Del(context.Background(), "refresh:"+claims.Username).Err()
		}
	}

	span.SetAttributes(attribute.String("logout.time", time.Now().String()))

	http.SetCookie(w, &http.Cookie{
		Name:     "Token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response := map[string]string{"message": "Logout Succes"}
	helpers.ResponseJSON(w, http.StatusOK, response)
}
