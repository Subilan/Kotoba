package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type reqCommonLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type reqCommonRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
	Website  string `json:"website"`
}

type reqCreateComment struct {
	Text string `json:"text"`
}

type reqDeleteComment struct {
	CommentId string `json:"comment_id"`
}

type reqReaction struct {
	CommentId string `json:"comment_id"`
	Emoji string `json:"emoji"`
}

type bsonComment struct {
	Username  string             `bson:"username"`
	Text      string             `bson:"text"`
	Uid       string             `bson:"uid"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}

type bsonCommentReaction struct {
	Username  string             `bson:"username"`
	Emoji     string             `bson:"emoji"`
	CommentId string             `bson:"comment_id"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}
