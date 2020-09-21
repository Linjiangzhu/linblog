package repository

import (
	"errors"
	"github.com/Linjiangzhu/linblog/linblog-backend/model"
	"gorm.io/gorm"
)

func (r *Repository) GetPost(pid uint) (*model.Post, error) {
	p := model.Post{ID: pid}
	err := r.db.First(&p).Error
	if err != nil {
		return nil, err
	}
	tags, _ := r.GetTagsByPID(pid)
	cats, _ := r.GetCatsByPID(pid)
	p.Tags = tags
	p.Categories = cats
	return &p, nil
}

func (r *Repository) GetVisiblePost(pid uint) (*model.Post, error) {
	p := model.Post{ID: pid}
	err := r.db.Where("visible = ?", true).First(&p).Error
	if err != nil {
		return nil, err
	}
	tags, _ := r.GetTagsByPID(pid)
	cats, _ := r.GetCatsByPID(pid)
	p.Tags = tags
	p.Categories = cats
	return &p, nil
}

func (r *Repository) GetTagsByPID(pid uint) ([]*model.Tag, error) {
	p := model.Post{ID: pid}
	var tags []model.Tag
	err := r.db.Model(&p).Association("Tags").Find(&tags)
	//err := r.db.Model(&p).Related(&tags, "Tags").Error
	if err != nil {
		return nil, err
	}
	var res []*model.Tag
	for _, t := range tags {
		res = append(res, &model.Tag{
			ID:   t.ID,
			Name: t.Name,
		})
	}
	return res, nil
}

func (r *Repository) GetCatsByPID(pid uint) ([]*model.Category, error) {
	p := model.Post{ID: pid}
	var cats []model.Category
	err := r.db.Model(&p).Association("Categories").Find(&cats)
	//err := r.db.Model(&p).Related(&cats, "Categories").Error
	if err != nil {
		return nil, err
	}
	var res []*model.Category
	for _, c := range cats {
		res = append(res, &model.Category{
			ID:   c.ID,
			Name: c.Name,
		})
	}
	return res, nil
}

func (r *Repository) CreatePost(p *model.Post) (res *model.Post, errs []error) {
	tags, terr := r.CreateTags(p.Tags)
	if len(terr) > 0 {
		errs = append(errs, terr...)
	}
	cats, cerr := r.CreateCats(p.Categories)
	if len(cerr) > 0 {
		errs = append(errs, cerr...)
	}
	p.Tags = tags
	p.Categories = cats
	err := r.db.Create(p).Error
	if err != nil {
		errs = append(errs, err)
	}
	res = p
	return
}

func (r *Repository) CreateTagIfNotExist(tag *model.Tag) (*model.Tag, error) {
	if err := r.db.Where("name = ?", tag.Name).First(tag).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		err := r.db.Create(tag).Error
		if err != nil {
			return nil, err
		}
	}
	return tag, nil
}

func (r *Repository) CreateCatIfNotExist(cat *model.Category) (*model.Category, error) {
	if err := r.db.Where("name = ?", cat.Name).First(cat).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		err := r.db.Create(cat).Error
		if err != nil {
			return nil, err
		}
	}
	return cat, nil
}

func (r *Repository) CreateTags(tags []*model.Tag) (res []*model.Tag, errs []error) {
	for _, tagPtr := range tags {
		tag, err := r.CreateTagIfNotExist(tagPtr)
		if err != nil {
			errs = append(errs, err)
		} else {
			res = append(res, tag)
		}
	}
	return
}

func (r *Repository) CreateCats(cats []*model.Category) (res []*model.Category, errs []error) {
	for _, catPtr := range cats {
		cat, err := r.CreateCatIfNotExist(catPtr)
		if err != nil {
			errs = append(errs, err)
		} else {
			res = append(res, cat)
		}
	}
	return
}

func (r *Repository) DeletePost(pid uint) (errs []error) {
	p := model.Post{ID: pid}
	err := r.db.First(&p).Error
	if err != nil {
		return []error{err}
	}
	err = r.db.Model(p).Association("Tags").Clear()
	if err != nil {
		errs = append(errs, err)
	}
	err = r.db.Model(p).Association("Categories").Clear()
	if err != nil {
		errs = append(errs, err)
	}
	err = r.db.Delete(p).Error
	if err != nil {
		errs = append(errs, err)
	}
	return
}

func (r *Repository) UpdatePost(p *model.Post) (res *model.Post, errs []error) {
	tags, terr := r.CreateTags(p.Tags)
	if len(terr) > 0 {
		errs = append(errs, terr...)
	}
	cats, cerr := r.CreateCats(p.Categories)
	if len(cerr) > 0 {
		errs = append(errs, cerr...)
	}
	err := r.db.Model(p).Association("Tags").Replace(tags)
	if err != nil {
		errs = append(errs, err)
	}
	err = r.db.Model(p).Association("Categories").Replace(cats)
	if err != nil {
		errs = append(errs, err)
	}
	err = r.db.Model(p).Update("visible", p.Visible).Error
	if err != nil {
		errs = append(errs, err)
	}
	err = r.db.Model(p).Updates(*p).Error
	if err != nil {
		errs = append(errs, err)
	}
	res = p
	r.db.First(res)
	return
}

func (r *Repository) GetPosts(offset, limit int) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Table("posts").Offset(offset).Limit(limit).Select("id, created_at, updated_at, " +
		"title, brief, cover, visible, user_id").Scan(&posts).Error
	if err != nil {
		return nil, err
	}
	for idx, p := range posts {
		tags, _ := r.GetTagsByPID(p.ID)
		cats, _ := r.GetCatsByPID(p.ID)
		posts[idx].Tags = tags
		posts[idx].Categories = cats
	}
	return posts, nil
}

func (r *Repository) GetVisiblePosts(offset, limit int) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Table("posts").Offset(offset).Limit(limit).Select("id, created_at, updated_at, "+
		"title, brief, cover, visible, user_id").Where("visible = ?", true).Scan(&posts).Error
	if err != nil {
		return nil, err
	}
	for idx, p := range posts {
		tags, _ := r.GetTagsByPID(p.ID)
		cats, _ := r.GetCatsByPID(p.ID)
		posts[idx].Tags = tags
		posts[idx].Categories = cats
	}
	return posts, nil
}
