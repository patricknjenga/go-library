package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"reflect"
)

type Fiber struct {
	*fiber.App
}

type query struct {
	Limit int `form:"limit"`
	Page  int `form:"page"`
}

func (f Fiber) Resource(p Postgres, model func() (any, any)) {
	_, r := model()
	route := reflect.TypeOf(r).Elem().Name()
	f.Get(fmt.Sprintf("/%s", route), func(c *fiber.Ctx) error {
		var count int64
		m, _ := model()
		var q = query{Page: 1, Limit: 10}
		err := c.QueryParser(&q)
		if err != nil {
			return err
		}
		err = p.Model(m).Count(&count).Error
		if err != nil {
			return err
		}
		err = p.Limit(q.Limit).Offset((q.Page - 1) * q.Limit).Find(m).Error
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{"count": count, "data": m, "limit": q.Limit, "page": q.Page})
	})
	f.Post(fmt.Sprintf("/%s", route), func(c *fiber.Ctx) error {
		_, m := model()
		err := c.BodyParser(m)
		if err != nil {
			return err
		}
		err = p.Create(m).Error
		if err != nil {
			return err
		}
		return c.JSON(m)
	})
	f.Post(fmt.Sprintf("/%s/batch", route), func(c *fiber.Ctx) error {
		m, _ := model()
		err := c.BodyParser(m)
		if err != nil {
			return err
		}
		err = p.Create(m).Error
		if err != nil {
			return err
		}
		return c.JSON(m)
	})
	f.Get(fmt.Sprintf("/%s/:id", route), func(c *fiber.Ctx) error {
		_, m := model()
		err := p.First(m, c.Params("id")).Error
		if err != nil {
			return fiber.ErrNotFound
		}
		return c.JSON(m)
	})
	f.Put(fmt.Sprintf("/%s/:id", route), func(c *fiber.Ctx) error {
		_, m := model()
		err := p.First(m, c.Params("id")).Error
		if err != nil {
			return fiber.ErrNotFound
		}
		err = c.BodyParser(m)
		if err != nil {
			return err
		}
		err = p.Updates(m).Error
		if err != nil {
			return err
		}
		return c.JSON(m)
	})
	f.Delete(fmt.Sprintf("/%s/:id", route), func(c *fiber.Ctx) error {
		_, m := model()
		q := p.Delete(m, c.Params("id"))
		if q.Error != nil {
			return q.Error
		}
		if q.RowsAffected == 0 {
			return fiber.ErrNotFound
		}
		return c.JSON(m)
	})
}
