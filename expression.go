package filter

import (
	"fmt"

	"gorm.io/gorm/clause"
)

type IExact clause.Eq

func (i IExact) Build(builder clause.Builder) {
	builder.WriteQuoted(i.Column)
	builder.WriteString(" ILIKE ")
	builder.AddVar(builder, i.Value)
}

type Contains clause.Eq

func (c Contains) Build(builder clause.Builder) {
	builder.WriteQuoted(c.Column)
	builder.WriteString(" LIKE ")
	builder.AddVar(builder, fmt.Sprintf("%%%s%%", c.Value))
}

type IContains clause.Eq

func (i IContains) Build(builder clause.Builder) {
	builder.WriteQuoted(i.Column)
	builder.WriteString(" ILIKE ")
	builder.AddVar(builder, fmt.Sprintf("%%%s%%", i.Value))
}

type StartsWith clause.Eq

func (s StartsWith) Build(builder clause.Builder) {
	builder.WriteQuoted(s.Column)
	builder.WriteString(" LIKE ")
	builder.AddVar(builder, fmt.Sprintf("%%%s", s.Value))
}

type IStartsWith clause.Eq

func (i IStartsWith) Build(builder clause.Builder) {
	builder.WriteQuoted(i.Column)
	builder.WriteString(" ILIKE ")
	builder.AddVar(builder, fmt.Sprintf("%%%s", i.Value))
}

type EndsWith clause.Eq

func (e EndsWith) Build(builder clause.Builder) {
	builder.WriteQuoted(e.Column)
	builder.WriteString(" LIKE ")
	builder.AddVar(builder, fmt.Sprintf("%s%%", e.Value))
}

type IEndsWith clause.Eq

func (i IEndsWith) Build(builder clause.Builder) {
	builder.WriteQuoted(i.Column)
	builder.WriteString(" ILIKE ")
	builder.AddVar(builder, fmt.Sprintf("%s%%", i.Value))
}
