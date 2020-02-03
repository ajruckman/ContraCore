package main

import (
    "github.com/gin-gonic/gin"

    . "github.com/ajruckman/xlib"

    "github.com/ajruckman/ContraCore/internal/db/contradb"
)

func main() {
    r := gin.Default()

    r.GET("/hourly", func(c *gin.Context){
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

        hourly, err := contradb.GetHourly()
        if err != nil {
            c.JSON(500, gin.H{
                "error": err.Error(),
            })
        } else {
            c.JSON(200, hourly)
        }
    })

    err := r.Run()
    Err(err)
}
