package clientsdk

import (
	"feiyu.com/wx/clientsdk/baseinfo"
	"feiyu.com/wx/protobuf/wechat"
	"github.com/golang/protobuf/proto"
)

// CreateSnsDeleteItem 创建操作项：删除指定的朋友圈项
func CreateSnsDeleteItem(snsObjID string) *baseinfo.SnsObjectOpItem {
	retItem := &baseinfo.SnsObjectOpItem{}
	//id,_:=strconv.ParseUint(opItems[index].SnsObjID,0,64)
	retItem.SnsObjID = snsObjID
	retItem.OpType = baseinfo.MMSnsOpCodeDelete
	retItem.DataLen = 0
	retItem.Data = []byte{}

	return retItem
}

// CreateSnsSetPrivateItem 创建操作项：设置自己发布的朋友圈项为私有信息(仅自己可见)
func CreateSnsSetPrivateItem(snsObjID string) *baseinfo.SnsObjectOpItem {
	retItem := &baseinfo.SnsObjectOpItem{}
	retItem.SnsObjID = snsObjID
	retItem.OpType = baseinfo.MMSnsOpCodeSetPrivate
	retItem.DataLen = 0
	retItem.Data = []byte{}

	return retItem
}

// CreateSnsSetPublicItem 创建操作项：设置自己发布的朋友圈项为公开信息(所有人可见)
func CreateSnsSetPublicItem(snsObjID string) *baseinfo.SnsObjectOpItem {
	retItem := &baseinfo.SnsObjectOpItem{}
	retItem.SnsObjID = snsObjID
	retItem.OpType = baseinfo.MMSnsOpCodeSetPublic
	retItem.DataLen = 0
	retItem.Data = []byte{}

	return retItem
}

// CreateSnsDeleteCommentItem 创建操作项：删除指定朋友圈项的评论
func CreateSnsDeleteCommentItem(snsObjID string, commentID uint32) *baseinfo.SnsObjectOpItem {
	retItem := &baseinfo.SnsObjectOpItem{}
	retItem.SnsObjID = snsObjID
	retItem.OpType = baseinfo.MMSnsOpCodeDeleteComment

	// 删除评论数据
	var opDeleteComment wechat.SnsObjectOpDeleteComment
	opDeleteComment.CommentId = &commentID
	data, _ := proto.Marshal(&opDeleteComment)
	retItem.DataLen = uint32(len(data))
	retItem.Data = data
	return retItem
}

// CreateSnsUnLikeItem 创建操作项：取消点赞
func CreateSnsUnLikeItem(snsObjID string) *baseinfo.SnsObjectOpItem {
	retItem := &baseinfo.SnsObjectOpItem{}
	retItem.SnsObjID = snsObjID
	retItem.OpType = baseinfo.MMSnsOpCodeUnLike
	retItem.DataLen = 0
	retItem.Data = []byte{}

	return retItem
}
