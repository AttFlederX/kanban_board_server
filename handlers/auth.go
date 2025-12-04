package handlers

import (
	"context"
	"time"

	"github.com/AttFlederX/kanban_board_server/middleware"
	"github.com/AttFlederX/kanban_board_server/models"
	"github.com/AttFlederX/kanban_board_server/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/api/idtoken"
)

type GoogleSignInRequest struct {
	IDToken string `json:"id_token"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func GoogleSignIn(c *fiber.Ctx, jwtSecret string) error {
	var req GoogleSignInRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			jsonFieldError: errInvalidRequestBody,
		})
	}

	if req.IDToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			jsonFieldError: errIDTokenRequired,
		})
	}

	// Verify the Google ID token
	payload, err := idtoken.Validate(context.Background(), req.IDToken, "")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			jsonFieldError: errInvalidGoogleToken,
		})
	}

	// Extract user information from the token payload
	googleID := payload.Subject
	email := payload.Claims[claimEmail].(string)
	name := payload.Claims[claimName].(string)
	photoURL := ""
	if pic, ok := payload.Claims[claimPicture].(string); ok {
		photoURL = pic
	}

	// Check if user exists in database
	var user models.User
	err = services.FindOne(collectionUsers, bson.M{fieldGoogleID: googleID}, &user)

	if err != nil {
		// User doesn't exist, create new user
		user = models.User{
			GoogleID: googleID,
			Email:    email,
			Name:     name,
			PhotoURL: photoURL,
		}

		userID, err := services.InsertOne(collectionUsers, user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				jsonFieldError: errFailedCreateUser,
			})
		}
		user.ID = userID
	} else {
		// User exists, update their information
		update := bson.M{
			fieldName:     name,
			fieldPhotoURL: photoURL,
			fieldEmail:    email,
		}
		if err := services.UpdateByID(collectionUsers, user.ID, update); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				jsonFieldError: errFailedUpdateUser,
			})
		}
		user.Name = name
		user.PhotoURL = photoURL
		user.Email = email
	}

	// Generate JWT token
	claims := middleware.Claims{
		UserID:   user.ID.Hex(),
		Email:    user.Email,
		GoogleID: user.GoogleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			jsonFieldError: errFailedGenerateToken,
		})
	}

	return c.JSON(AuthResponse{
		Token: tokenString,
		User:  user,
	})
}
