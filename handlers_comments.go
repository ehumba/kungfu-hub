package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ehumba/kungfu-hub/internal/auth"
	"github.com/ehumba/kungfu-hub/internal/database"
	"github.com/google/uuid"
)

func (a *apiConfig) postCommentHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "no authorization header")
		return
	}

	userID, err := auth.ValidateJWT(token, a.secret)
	if err != nil {
		respondWithError(w, 401, "invalid token")
		return
	}

	var params struct {
		NewsItemID *int32 `json:"news_item_id,omitempty"`
		ParentID   *int32 `json:"parent_id,omitempty"`
		Content    string `json:"content"`
	}

	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil || params.Content == "" {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if params.NewsItemID == nil && params.ParentID == nil {
		respondWithError(w, http.StatusBadRequest, "either news_item_id or parent_id must be provided")
		return
	}

	var news_item_id sql.NullInt32
	var parent_id sql.NullInt32

	if params.NewsItemID == nil && params.ParentID != nil {
		parent_id = sql.NullInt32{Int32: *params.ParentID, Valid: true}
		// Fetch the parent comment to get its NewsItemID
		parentComment, err := a.dbQueries.GetCommentByID(r.Context(), parent_id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "parent comment not found")
			return
		}
		if parentComment.NewsItemID.Valid == false {
			respondWithError(w, http.StatusBadRequest, "cannot reply to a reply")
			return
		}
		news_item_id = parentComment.NewsItemID
	}

	if params.NewsItemID != nil && params.ParentID == nil {
		news_item_id = sql.NullInt32{Int32: *params.NewsItemID, Valid: true}
		parent_id = sql.NullInt32{Valid: false}
	}

	if params.NewsItemID != nil && params.ParentID != nil {
		respondWithError(w, http.StatusBadRequest, "provide either news_item_id or parent_id, not both")
		return
	}

	var userIDNull uuid.NullUUID
	userIDNull.UUID = userID
	userIDNull.Valid = true

	var content sql.NullString
	content.String = params.Content
	content.Valid = true

	newCommentParams := database.PostCommentParams{
		Column1: news_item_id,
		Column2: userIDNull,
		Column3: parent_id,
		Column4: content,
	}

	newCommentDB, err := a.dbQueries.PostComment(r.Context(), newCommentParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not post comment")
		return
	}

	newComment := Comment{
		ID:              newCommentDB.ID,
		NewsItemID:      newCommentDB.NewsItemID,
		UserID:          newCommentDB.UserID,
		ParentCommentID: newCommentDB.ParentCommentID,
		Content:         newCommentDB.Content,
		CreatedAt:       newCommentDB.CreatedAt,
		UpdatedAt:       newCommentDB.UpdatedAt,
	}

	respondWithJSON(w, http.StatusCreated, newComment)
}

func (a *apiConfig) editCommentHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "no authorization header")
		return
	}

	userID, err := auth.ValidateJWT(token, a.secret)
	if err != nil {
		respondWithError(w, 401, "invalid token")
		return
	}

	var params struct {
		CommentID int32  `json:"comment_id"`
		Content   string `json:"content"`
	}

	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil || params.Content == "" {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	id := sql.NullInt32{Int32: params.CommentID, Valid: true}
	content := sql.NullString{String: params.Content, Valid: true}
	var userIDNull uuid.NullUUID
	userIDNull.UUID = userID
	userIDNull.Valid = true

	// Fetch the comment to verify ownership
	comment, err := a.dbQueries.GetCommentByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "comment not found")
		return
	}

	if comment.UserID.Valid == false || comment.UserID.UUID != userID {
		respondWithError(w, http.StatusForbidden, "you do not have permission to edit this comment")
		return
	}

	editCommentParams := database.EditCommentParams{
		Column1: id,
		Column2: content,
	}
	editedCommentDB, err := a.dbQueries.EditComment(r.Context(), editCommentParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not edit comment")
		return
	}

	editedComment := Comment{
		ID:              editedCommentDB.ID,
		NewsItemID:      editedCommentDB.NewsItemID,
		UserID:          editedCommentDB.UserID,
		ParentCommentID: editedCommentDB.ParentCommentID,
		Content:         editedCommentDB.Content,
		CreatedAt:       editedCommentDB.CreatedAt,
		UpdatedAt:       editedCommentDB.UpdatedAt,
	}

	respondWithJSON(w, http.StatusOK, editedComment)
}

func (a *apiConfig) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "no authorization header")
		return
	}

	userID, err := auth.ValidateJWT(token, a.secret)
	if err != nil {
		respondWithError(w, 401, "invalid token")
		return
	}

	var params struct {
		CommentID int32 `json:"comment_id"`
	}

	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	id := sql.NullInt32{Int32: params.CommentID, Valid: true}

	var userIDNull uuid.NullUUID
	userIDNull.UUID = userID
	userIDNull.Valid = true

	// Fetch the comment to verify ownership
	comment, err := a.dbQueries.GetCommentByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "comment not found")
		return
	}

	if comment.UserID.Valid == false || comment.UserID.UUID != userID {
		respondWithError(w, http.StatusForbidden, "you do not have permission to delete this comment")
		return
	}

	err = a.dbQueries.DeleteComment(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete comment")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
