package core

import (
	"log"
)

type Change struct {
	New map[int][]rune
	Old map[int][]rune
}

func Insert(y int, r *Row) Change {
	return Change{
		New: map[int][]rune{},
		Old: map[int][]rune{},
	}
}

func Edit(y int, r []rune) Change {
	return Change{
		New: map[int][]rune{y: r},
		Old: map[int][]rune{},
	}
}

func Delete(from, to int) Change {
	rows := make(map[int][]rune, to-from)

	for i := from; i <= to; i++ {
		rows[i] = nil
	}

	return Change{New: rows, Old: make(map[int][]rune)}
}

func (c Change) Apply(e *E) {
	doFullRender := false

	for i, row := range c.New {
		log.Printf("applying change i=%d, row=%+v", i, row)

		// edit needs to store the previous row data
		if _, ok := c.Old[i]; !ok {
			c.Old[i] = make([]rune, len(e.Row(i)))
			copy(c.Old[i], e.Row(i))

			e.SetRow(i, row)

			if row != nil {
				e.Render(i)
			} else {
				doFullRender = true
			}
		} else {
			e.rows = append(append(e.rows[:i], &Row{chars:row}), e.rows[i:]...)
			doFullRender = true
		}
	}

	Filter(&e.rows, func(i int) bool { return e.rows[i] != nil })

	if doFullRender {
		e.FullRender()
	}
}

func (c Change) Undo() Change {
	return Change{New: c.Old, Old: c.New}
}
