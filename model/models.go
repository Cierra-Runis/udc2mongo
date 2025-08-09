package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// See: https://www.unicode.org/reports/tr42/#d1e2729
type UCD struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	Xmlns       string      `xml:"xmlns,attr" bson:"-" json:"xmlns,omitempty"`      // See: https://www.unicode.org/reports/tr42/#d1e2673
	Description string      `xml:"description" json:"description,omitempty"`        // See: https://www.unicode.org/reports/tr42/#d1e2800
	Repertoire  *Repertoire `xml:"repertoire" bson:"-" json:"repertoire,omitempty"` // See: https://www.unicode.org/reports/tr42/#d1e2832
	Blocks      *Blocks     `xml:"blocks" bson:"-" json:"blocks,omitempty"`         // See: https://www.unicode.org/reports/tr42/#d1e3971

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	Version   string    `bson:"version" json:"version"`
}

// Blocks 字符块定义
type Blocks struct {
	Blocks []Block `xml:"block" bson:"blocks" json:"blocks,omitempty"`
}

// Block 字符块
type Block struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstCP string             `xml:"first-cp,attr" bson:"first_cp" json:"first_cp"`
	LastCP  string             `xml:"last-cp,attr" bson:"last_cp" json:"last_cp"`
	Name    string             `xml:"name,attr" bson:"name" json:"name"`

	// 时间戳
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
