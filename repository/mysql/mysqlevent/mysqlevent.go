package mysqlevent

import (
	"database/sql"
	"event-manager/entity"
	"event-manager/repository/mysql"
	"fmt"
	"log"
	"time"
)

type EventRepo struct {
	conn *mysql.MySQLDB
}

func New(c *mysql.MySQLDB) EventRepo {
	return EventRepo{
		conn: c,
	}
}

type eventModel struct {
	id        uint
	ownerID   uint
	title     string
	location  string
	startAt   time.Time
	status    entity.EventStatus
	createdAt time.Time
}

func (e *eventModel) ToEventEntity() entity.Event {
	var entiy entity.Event

	entiy.SetID(e.id)
	entiy.SetOwner(e.ownerID)
	entiy.SetTitle(e.title)
	entiy.SetLocation(e.location)
	entiy.SetStartAt(e.startAt)
	entiy.SetStatus(e.status)

	return entiy
}

func (r EventRepo) GetEventByID(id uint) (entity.Event, bool, error) {
	var model eventModel

	row := r.conn.Conn().QueryRow(`select id, title from events where id = ?`, id)
	err := row.Scan(&model.id, &model.title)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Event{}, false, nil
		}

		return entity.Event{}, false, err
	}

	return model.ToEventEntity(), true, nil
}

func (r EventRepo) GetAllEventsFor(userId uint) ([]entity.Event, error) {
	eventACL := make([]entity.Event, 0)

	rows, err := r.conn.Conn().Query(`select id, title, location, start_at from events where owner_id = ?`, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var acl eventModel

		err := rows.Scan(&acl.id, &acl.title, &acl.location, &acl.startAt)
		if err != nil {
			return nil, err
		}

		eventACL = append(eventACL, acl.ToEventEntity())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return eventACL, nil
}

func (r EventRepo) CreateEvent(e entity.Event) (entity.Event, error) {

	res, err := r.conn.Conn().Exec("INSERT INTO events (title, location, start_at, owner_id, status, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		e.Title(), e.Location(), e.StartAt(), e.OwnerID(), e.Status(), time.Now())

	if err != nil {
		return e, err
	}

	id, _ := res.LastInsertId()
	e.SetID(uint(id))
	
	return e, nil
}

func (r EventRepo) UpdateEvent(event entity.Event) error {
	_, err := r.conn.Conn().Exec("UPDATE events SET title=?, location=? WHERE id=?", event.Title(), event.Location(), event.ID())
	if err != nil {
		return err
	}

	return nil
}

func (r EventRepo) DeleteEvent(id uint) error {
	log.Print("im in repo layer... \n")
	_, err := r.conn.Conn().Exec(`DELETE FROM events WHERE id = ?`, id)

	if err != nil {
		return fmt.Errorf("can't execute command: %w", err)
	}

	return nil
}
