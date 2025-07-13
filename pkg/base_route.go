package pkg

import (
	"encoding/json"
	"reflect"

	"github.com/gin-gonic/gin"
)

type RouteOptions struct {
	EnableList bool
	ListMsgID  string

	EnableGet bool
	GetMsgID  string

	EnableCreate bool
	CreateMsgID  string

	EnableUpdate bool
	UpdateMsgID  string

	EnableDelete bool
	DeleteMsgID  string
}

type CRUDHandlers struct {
	List   gin.HandlerFunc
	Get    gin.HandlerFunc
	Create gin.HandlerFunc
	Update gin.HandlerFunc
	Delete gin.HandlerFunc
}

type RequestPayload[T any] struct {
	Header struct {
		MsgID string `json:"msgId"`
	} `json:"header"`
	Data T `json:"data"`
}

type IDPayload struct {
	ID uint `json:"id" binding:"required"`
}

type RequestTypes struct {
	List   any
	Get    any
	Create any
	Update any
	Delete any
}

// TODO: problem is not check with msgId

// func RegisterUnifiedRoute[T any](
// 	group *gin.RouterGroup,
// 	svc *BaseService[T],
// 	opts RouteOptions,
// 	uniqueFields []string,
// 	handlers *CRUDHandlers,
// 	requestTypes *RequestTypes,
// ) {
// 	handlerMap := map[string]gin.HandlerFunc{}

// 	if opts.EnableList {
// 		var req any
// 		if requestTypes != nil {
// 			req = requestTypes.List
// 		}
// 		if handlers != nil && handlers.List != nil {
// 			handlerMap[opts.ListMsgID] = handlers.List
// 		} else {
// 			handlerMap[opts.ListMsgID] = wrapListHandler(svc, req)
// 		}
// 	}

// 	if opts.EnableGet {
// 		var req any
// 		if requestTypes != nil {
// 			req = requestTypes.Get
// 		}
// 		if handlers != nil && handlers.Get != nil {
// 			handlerMap[opts.GetMsgID] = handlers.Get
// 		} else {
// 			handlerMap[opts.GetMsgID] = wrapGetHandler(svc, req)
// 		}
// 	}

// 	if opts.EnableCreate {
// 		var req any
// 		if requestTypes != nil {
// 			req = requestTypes.Create
// 		}
// 		if handlers != nil && handlers.Create != nil {
// 			handlerMap[opts.CreateMsgID] = handlers.Create
// 		} else {
// 			handlerMap[opts.CreateMsgID] = wrapCreateHandler(svc, uniqueFields, req)
// 		}
// 	}

// 	if opts.EnableUpdate {
// 		var req any
// 		if requestTypes != nil {
// 			req = requestTypes.Update
// 		}
// 		if handlers != nil && handlers.Update != nil {
// 			handlerMap[opts.UpdateMsgID] = handlers.Update
// 		} else {
// 			handlerMap[opts.UpdateMsgID] = wrapUpdateHandler(svc, req)
// 		}
// 	}

// 	if opts.EnableDelete {
// 		var req any
// 		if requestTypes != nil {
// 			req = requestTypes.Delete
// 		}
// 		if handlers != nil && handlers.Delete != nil {
// 			handlerMap[opts.DeleteMsgID] = handlers.Delete
// 		} else {
// 			handlerMap[opts.DeleteMsgID] = wrapDeleteHandler(svc, req)
// 		}
// 	}

// 	group.POST("/doProcess", func(c *gin.Context) {
// 		var raw RequestPayload[map[string]interface{}]
// 		if err := c.ShouldBindJSON(&raw); err != nil {
// 			ErrorWithStatus(c, "invalid payload")
// 			return
// 		}
// 		handler, ok := handlerMap[raw.Header.MsgID]
// 		if !ok {
// 			ErrorWithStatus(c, "unknown msgId")
// 			return
// 		}
// 		c.Set("rawPayload", raw)
// 		handler(c)
// 	})
// }

func RegisterUnifiedRoute[T any](
	router *gin.Engine,
	svc *BaseService[T],
	opts RouteOptions,
	uniqueFields []string,
	handlers *CRUDHandlers,
	requestTypes *RequestTypes,
) {
	router.POST("/api/doProcess", func(c *gin.Context) {
		var payload RequestPayload[json.RawMessage]
		if err := c.ShouldBindJSON(&payload); err != nil {
			ErrorWithStatus(c, "invalid request format")
			return
		}

		switch payload.Header.MsgID {
		case opts.ListMsgID:
			(wrapListHandler(svc, requestTypes.List))(c)
		case opts.GetMsgID:
			(wrapGetHandler(svc, requestTypes.Get))(c)
		case opts.CreateMsgID:
			(wrapCreateHandler(svc, uniqueFields, requestTypes.Create))(c)
		case opts.UpdateMsgID:
			(wrapUpdateHandler(svc, requestTypes.Update))(c)
		case opts.DeleteMsgID:
			(wrapDeleteHandler(svc, requestTypes.Delete))(c)
		default:
			ErrorWithStatus(c, "unknown msgId")
		}
	})
}


// ======================= WRAPPED HANDLERS ==============================

func wrapCreateHandler[T any](svc *BaseService[T], uniqueFields []string, customType any) gin.HandlerFunc {
	return func(c *gin.Context) {
		if customType == nil {
			var req RequestPayload[T]
			if err := c.ShouldBindJSON(&req); err != nil {
				ErrorWithStatus(c, "invalid data")
				return
			}
			for _, field := range uniqueFields {
				exists, err := svc.Exists(field, req.Data)
				if err != nil {
					ErrorWithStatus(c, "validation failed")
					return
				}
				if exists {
					ErrorWithStatus(c, "duplicate "+field)
					return
				}
			}
			if err := svc.Create(&req.Data); err != nil {
				ErrorWithStatus(c, "failed to create")
				return
			}
			Success(c, "created successfully", req.Data)
			return
		}

		var raw RequestPayload[map[string]interface{}]
		if err := c.ShouldBindJSON(&raw); err != nil {
			ErrorWithStatus(c, "invalid data")
			return
		}

		rawBytes, _ := json.Marshal(raw.Data)
		typedData := reflect.New(reflect.TypeOf(new(T)).Elem()).Interface()
		_ = json.Unmarshal(rawBytes, typedData)

		data := typedData.(*T)

		for _, field := range uniqueFields {
			exists, err := svc.Exists(field, *data)
			if err != nil {
				ErrorWithStatus(c, "validation failed")
				return
			}
			if exists {
				ErrorWithStatus(c, "duplicate "+field)
				return
			}
		}

		if err := svc.Create(data); err != nil {
			ErrorWithStatus(c, "failed to create")
			return
		}

		Success(c, "created successfully", data)
	}
}

func wrapListHandler[T any](svc *BaseService[T], _ any) gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := svc.List()
		if err != nil {
			ErrorWithStatus(c, "failed to list data")
			return
		}
		Success(c, "", items)
	}
}

func wrapGetHandler[T any](svc *BaseService[T], _ any) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := bindCustomPayload[IDPayload](c, nil)
		if err != nil {
			ErrorWithStatus(c, "invalid ID payload")
			return
		}
		item, err := svc.GetByID(payload.Data.ID)
		if err != nil {
			ErrorWithStatus(c, "data not found")
			return
		}
		Success(c, "", item)
	}
}

func wrapUpdateHandler[T any](svc *BaseService[T], reqType any) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := bindCustomPayload[map[string]interface{}](c, reqType)
		if err != nil {
			ErrorWithStatus(c, "invalid update data")
			return
		}

		rawID, ok := payload.Data["id"]
		if !ok {
			ErrorWithStatus(c, "missing id")
			return
		}
		idFloat, ok := rawID.(float64)
		if !ok {
			ErrorWithStatus(c, "invalid id format")
			return
		}
		id := uint(idFloat)

		// 🚨 Remove 'id' field to avoid trying to update the primary key
		delete(payload.Data, "id")

		if err := svc.Update(id, payload.Data); err != nil {
			ErrorWithStatus(c, err.Error())
			return
		}
		Success(c, "updated successfully", nil)
	}
}

func wrapDeleteHandler[T any](svc *BaseService[T], _ any) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := bindCustomPayload[IDPayload](c, nil)
		if err != nil {
			ErrorWithStatus(c, "invalid ID")
			return
		}
		if err := svc.Delete(payload.Data.ID); err != nil {
			ErrorWithStatus(c, err.Error())
			return
		}
		Success(c, "deleted successfully", nil)
	}
}

// ======================= BINDING UTILITY ==============================

func bindCustomPayload[T any](c *gin.Context, _ any) (*RequestPayload[T], error) {
	var payload RequestPayload[T]
	if err := c.ShouldBindJSON(&payload); err != nil {
		return nil, err
	}
	return &payload, nil
}
