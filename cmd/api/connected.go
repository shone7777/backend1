package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/nimilgp/URLcommentary/internal/dblayer"
	"github.com/nimilgp/URLcommentary/internal/graphdb"
)

func (app *application) getConnectedComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	pageID, err := strconv.Atoi(ps.ByName("pageid"))
	if err != nil {
		http.Error(w, "Invalid page ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	const LIMIT int32 = 10

	comments, err := app.queries.GetCommentsByPage(ctx, dblayer.GetCommentsByPageParams{
		Pageid: int32(pageID),
		Limit:  LIMIT,
	})
	if err != nil {
		http.Error(w, "Error fetching comments", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Fetched %d comments:\n", len(comments))
	for i, comment := range comments {
		fmt.Printf("[%d] CommentID: %d, UserID: %d, Comment: %s, CreatedAt: %v\n",
			i, comment.Commentid, comment.Userid, comment.Commentdata, comment.Createdat)
	}

	type CommentWithHops struct {
		CommentID int32  `json:"comment_id"`
		UserID    int32  `json:"user_id"`
		Comment   string `json:"comment"`
		CreatedAt string `json:"created_at"`
		HopCount  int    `json:"hop_count"`
	}

	var commentsWithHops []CommentWithHops
	fmt.Println("Entering getConnectedComments")

	for _, comment := range comments {
		hopCount, err := graphdb.GetHopCount(int64(userID), int64(comment.Userid))
		fmt.Printf("Calling GetHopCount with UserID: %d, Comment UserID: %d\n", userID, comment.Userid)
		if err != nil {
			fmt.Printf("Error in GetHopCount: %v\n", err)
			hopCount = -1
		} else {
			fmt.Printf("HopCount result: %d\n", hopCount)
		}

		commentsWithHops = append(commentsWithHops, CommentWithHops{
			CommentID: comment.Commentid,
			UserID:    int32(comment.Userid),
			Comment:   comment.Commentdata,
			CreatedAt: comment.Createdat.Time.Format("2006-01-02 15:04:05"),
			HopCount:  hopCount,
		})
	}

	sort.Slice(commentsWithHops, func(i, j int) bool {
		if commentsWithHops[i].HopCount == -1 {
			return false
		}
		if commentsWithHops[j].HopCount == -1 {
			return true
		}
		return commentsWithHops[i].HopCount < commentsWithHops[j].HopCount
	})

	app.writeJSON(w, http.StatusOK, envelope{"comments": commentsWithHops}, nil)
}

func (app *application) handleConnect(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	userid, err := strconv.Atoi(ps.ByName("userid"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	commentid, err := strconv.Atoi(ps.ByName("commentid"))
	if err != nil {
		http.Error(w, "Invalid commentID", http.StatusBadRequest)
		return
	}
	mode, err := strconv.Atoi(ps.ByName("mode"))
	if err != nil {
		http.Error(w, "Invalid mode", http.StatusBadRequest)
		return
	}
	user2ID, err := app.queries.GetUserIDByCommentID(ctx, int32(commentid))
	if err != nil {
		http.Error(w, "Error fetching userid", http.StatusInternalServerError)
		return
	}
	if mode == 1 {
		err = graphdb.ConnectUsers(int64(userid), int64(user2ID))
	} else {
		err = graphdb.DisconnectUsers(int64(userid), int64(user2ID))
	}
	if err != nil {
		http.Error(w, "Failed to process request", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
