package api

import (
	"github.com/gin-gonic/gin"
	"mngr/eb"
	"mngr/models"
	"mngr/utils"
	"net/http"
	"path"
	"time"
)

func RegisterSourceEndpoints(router *gin.Engine) {
	router.GET("/sources", func(ctx *gin.Context) {
		sources, _ := utils.SourceRep.GetAll()
		ctx.JSON(http.StatusOK, sources)
	})
	router.GET("/sources/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		source, err := utils.SourceRep.Get(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, source)
	})
	router.POST("/sources", func(ctx *gin.Context) {
		var model models.SourceModel
		if err := ctx.ShouldBindJSON(&model); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		isNew := len(model.Id) == 0
		model.CreatedAt = utils.FromDateToString(time.Now())
		if _, err := utils.SourceRep.Save(&model); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//restart the stream after edit.
		if !isNew {
			eventPub := eb.RestartStreamRequestEvent{SourceModel: model}
			err := eventPub.Publish()
			if err == nil {
				ctx.Writer.WriteHeader(http.StatusOK)
			}
		}
		// Create stream folder.
		sf, _ := utils.GetStreamFolderPath()
		dic := path.Join(sf, model.Id)
		utils.CreateDicIfNotExist(dic)

		// Create record folder.
		sf, _ = utils.GetRecordFolderPath()
		dic = path.Join(sf, model.Id)
		utils.CreateDicIfNotExist(dic)

		// Create read folder.
		sf, _ = utils.GetReadFolderPath()
		dic = path.Join(sf, model.Id)
		utils.CreateDicIfNotExist(dic)

		ctx.JSON(http.StatusOK, model)
	})
	router.DELETE("/sources/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if err := utils.SourceRep.RemoveById(id); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//stops the stream after delete.
		ssrEvent := eb.StopStreamRequestEvent{Id: id}
		err := ssrEvent.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}

		//also remove Object Detection Model
		if err := utils.OdRep.RemoveById(id); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//

		mc := eb.ModelChanged{SourceId: id}
		mcJson, _ := utils.SerializeJson(mc)
		dcEvent := eb.DataChangedEvent{ModelName: "od", ParamsJson: mcJson, Op: eb.DELETE}
		err = dcEvent.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
		ctx.JSON(http.StatusOK, gin.H{"id": id})
	})
	router.GET("/sourcestreamstatus", func(context *gin.Context) {
		modelList, err := utils.SourceRep.GetSourceStreamStatus(utils.StreamRep)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusOK, modelList)
	})
}
