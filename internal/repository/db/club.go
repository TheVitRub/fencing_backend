package db

import (
	"context"
	"fmt"
	"time"

	dbEntities "fencing-club/internal/domain/entities/db"

	"github.com/lib/pq"
)

func (r *DBRepository) ListComments(ctx context.Context, targetType string, targetID int64, includeHidden bool) ([]dbEntities.Comment, error) {
	var comments []dbEntities.Comment
	statusWhere := `AND c.status='visible'`
	if includeHidden {
		statusWhere = ``
	}
	err := r.db.Select(ctx, nil, &comments,
		`SELECT c.id, c.target_type, c.target_id, c.user_id,
		        COALESCE(NULLIF(u.display_name, ''), u.login) AS user_display_name,
		        u.role AS user_role,
		        c.body, c.status, c.created_at, c.updated_at
		 FROM comments c
		 JOIN users u ON u.id=c.user_id
		 WHERE c.target_type=$1 AND c.target_id=$2 `+statusWhere+`
		 ORDER BY c.created_at ASC`,
		targetType, targetID,
	)
	if err != nil {
		return nil, fmt.Errorf("ListComments: %w", err)
	}
	return comments, nil
}

func (r *DBRepository) CreateComment(ctx context.Context, c *dbEntities.Comment) error {
	err := r.db.QueryRow(ctx, nil,
		[]any{&c.ID},
		`INSERT INTO comments (target_type, target_id, user_id, body, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, 'visible', $5, $5)
		 RETURNING id`,
		c.TargetType, c.TargetID, c.UserID, c.Body, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("CreateComment: %w", err)
	}
	return nil
}

func (r *DBRepository) UpdateCommentStatus(ctx context.Context, id int64, status string) error {
	_, err := r.db.Exec(ctx, nil,
		`UPDATE comments SET status=$1, updated_at=NOW() WHERE id=$2`,
		status, id,
	)
	if err != nil {
		return fmt.Errorf("UpdateCommentStatus id=%d: %w", id, err)
	}
	return nil
}

func (r *DBRepository) DeleteComment(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, nil, `UPDATE comments SET status='deleted', updated_at=NOW() WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("DeleteComment id=%d: %w", id, err)
	}
	return nil
}

func (r *DBRepository) SetEventAttendance(ctx context.Context, eventID, userID int64, status string) error {
	_, err := r.db.Exec(ctx, nil,
		`INSERT INTO event_attendees (event_id, user_id, status)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (event_id, user_id)
		 DO UPDATE SET status=EXCLUDED.status, created_at=NOW()`,
		eventID, userID, status,
	)
	if err != nil {
		return fmt.Errorf("SetEventAttendance: %w", err)
	}
	return nil
}

func (r *DBRepository) ListEventAttendees(ctx context.Context, eventID int64) ([]dbEntities.EventAttendee, error) {
	var attendees []dbEntities.EventAttendee
	err := r.db.Select(ctx, nil, &attendees,
		`SELECT ea.event_id, ea.user_id,
		        COALESCE(NULLIF(u.display_name, ''), u.login) AS user_display_name,
		        u.role AS user_role,
		        ea.status, ea.created_at
		 FROM event_attendees ea
		 JOIN users u ON u.id=ea.user_id
		 WHERE ea.event_id=$1 AND ea.status='going'
		 ORDER BY ea.created_at ASC`,
		eventID,
	)
	if err != nil {
		return nil, fmt.Errorf("ListEventAttendees: %w", err)
	}
	return attendees, nil
}

func (r *DBRepository) GetEventAttendance(ctx context.Context, eventID, userID int64) (*dbEntities.EventAttendee, error) {
	var a dbEntities.EventAttendee
	err := r.db.Get(ctx, &a,
		`SELECT ea.event_id, ea.user_id,
		        COALESCE(NULLIF(u.display_name, ''), u.login) AS user_display_name,
		        u.role AS user_role,
		        ea.status, ea.created_at
		 FROM event_attendees ea
		 JOIN users u ON u.id=ea.user_id
		 WHERE ea.event_id=$1 AND ea.user_id=$2`,
		eventID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("GetEventAttendance: %w", err)
	}
	return &a, nil
}

func (r *DBRepository) CreateNotificationsForEventAttendees(ctx context.Context, eventID int64, title, body string) error {
	_, err := r.db.Exec(ctx, nil,
		`INSERT INTO notifications (user_id, event_id, title, body)
		 SELECT user_id, event_id, $2, $3
		 FROM event_attendees
		 WHERE event_id=$1 AND status='going'`,
		eventID, title, body,
	)
	if err != nil {
		return fmt.Errorf("CreateNotificationsForEventAttendees: %w", err)
	}
	return nil
}

func (r *DBRepository) ListNotifications(ctx context.Context, userID int64) ([]dbEntities.Notification, error) {
	var notifications []dbEntities.Notification
	err := r.db.Select(ctx, nil, &notifications,
		`SELECT id, user_id, event_id, title, body, is_read, created_at
		 FROM notifications
		 WHERE user_id=$1
		 ORDER BY created_at DESC
		 LIMIT 50`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("ListNotifications: %w", err)
	}
	return notifications, nil
}

func (r *DBRepository) MarkNotificationRead(ctx context.Context, userID, notificationID int64) error {
	_, err := r.db.Exec(ctx, nil,
		`UPDATE notifications SET is_read=TRUE WHERE id=$1 AND user_id=$2`,
		notificationID, userID,
	)
	if err != nil {
		return fmt.Errorf("MarkNotificationRead: %w", err)
	}
	return nil
}

func (r *DBRepository) ListInstructorProfiles(ctx context.Context) ([]dbEntities.InstructorProfile, error) {
	var profiles []dbEntities.InstructorProfile
	err := r.db.Select(ctx, nil, &profiles,
		`SELECT p.id, p.user_id, u.role, p.name, p.photo_url, p.specialization,
		        p.weapons, p.experience, p.quote, p.bio, p.updated_at
		 FROM instructor_profiles p
		 JOIN users u ON u.id=p.user_id
		 ORDER BY p.name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("ListInstructorProfiles: %w", err)
	}
	return profiles, nil
}

func (r *DBRepository) UpsertInstructorProfile(ctx context.Context, p *dbEntities.InstructorProfile) error {
	err := r.db.QueryRow(ctx, nil,
		[]any{&p.ID},
		`INSERT INTO instructor_profiles
		   (user_id, name, photo_url, specialization, weapons, experience, quote, bio, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW())
		 ON CONFLICT (user_id) DO UPDATE SET
		   name=EXCLUDED.name, photo_url=EXCLUDED.photo_url, specialization=EXCLUDED.specialization,
		   weapons=EXCLUDED.weapons, experience=EXCLUDED.experience, quote=EXCLUDED.quote,
		   bio=EXCLUDED.bio, updated_at=NOW()
		 RETURNING id`,
		p.UserID, p.Name, p.PhotoURL, p.Specialization, p.Weapons, p.Experience, p.Quote, p.Bio,
	)
	if err != nil {
		return fmt.Errorf("UpsertInstructorProfile: %w", err)
	}
	return nil
}

func (r *DBRepository) DeleteInstructorProfile(ctx context.Context, userID int64) error {
	_, err := r.db.Exec(ctx, nil, `DELETE FROM instructor_profiles WHERE user_id=$1`, userID)
	if err != nil {
		return fmt.Errorf("DeleteInstructorProfile: %w", err)
	}
	return nil
}

func (r *DBRepository) ListKnowledgeArticles(ctx context.Context, role string) ([]dbEntities.KnowledgeArticle, error) {
	var articles []dbEntities.KnowledgeArticle
	allowed := []string{"public"}
	if role != "" {
		allowed = append(allowed, "registered")
	}
	if role == "student" || role == "instructor" || role == "admin" || role == "founder" {
		allowed = append(allowed, "student")
	}
	err := r.db.Select(ctx, nil, &articles,
		`SELECT id, title, category, body, visibility, sort_order, created_at, updated_at
		 FROM knowledge_articles
		 WHERE visibility = ANY($1)
		 ORDER BY sort_order ASC, title ASC`,
		pq.Array(allowed),
	)
	if err != nil {
		return nil, fmt.Errorf("ListKnowledgeArticles: %w", err)
	}
	return articles, nil
}

func (r *DBRepository) CreateKnowledgeArticle(ctx context.Context, a *dbEntities.KnowledgeArticle) error {
	err := r.db.QueryRow(ctx, nil,
		[]any{&a.ID},
		`INSERT INTO knowledge_articles (title, category, body, visibility, sort_order)
		 VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		a.Title, a.Category, a.Body, a.Visibility, a.SortOrder,
	)
	if err != nil {
		return fmt.Errorf("CreateKnowledgeArticle: %w", err)
	}
	return nil
}

func (r *DBRepository) UpdateKnowledgeArticle(ctx context.Context, a *dbEntities.KnowledgeArticle) error {
	_, err := r.db.Exec(ctx, nil,
		`UPDATE knowledge_articles
		 SET title=$1, category=$2, body=$3, visibility=$4, sort_order=$5, updated_at=NOW()
		 WHERE id=$6`,
		a.Title, a.Category, a.Body, a.Visibility, a.SortOrder, a.ID,
	)
	if err != nil {
		return fmt.Errorf("UpdateKnowledgeArticle: %w", err)
	}
	return nil
}

func (r *DBRepository) DeleteKnowledgeArticle(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, nil, `DELETE FROM knowledge_articles WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("DeleteKnowledgeArticle: %w", err)
	}
	return nil
}

func (r *DBRepository) ListGlossaryTerms(ctx context.Context) ([]dbEntities.GlossaryTerm, error) {
	var terms []dbEntities.GlossaryTerm
	err := r.db.Select(ctx, nil, &terms,
		`SELECT id, term, category, definition, created_at, updated_at
		 FROM glossary_terms
		 ORDER BY term ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("ListGlossaryTerms: %w", err)
	}
	return terms, nil
}

func (r *DBRepository) CreateGlossaryTerm(ctx context.Context, t *dbEntities.GlossaryTerm) error {
	err := r.db.QueryRow(ctx, nil,
		[]any{&t.ID},
		`INSERT INTO glossary_terms (term, category, definition)
		 VALUES ($1,$2,$3) RETURNING id`,
		t.Term, t.Category, t.Definition,
	)
	if err != nil {
		return fmt.Errorf("CreateGlossaryTerm: %w", err)
	}
	return nil
}

func (r *DBRepository) UpdateGlossaryTerm(ctx context.Context, t *dbEntities.GlossaryTerm) error {
	_, err := r.db.Exec(ctx, nil,
		`UPDATE glossary_terms SET term=$1, category=$2, definition=$3, updated_at=NOW() WHERE id=$4`,
		t.Term, t.Category, t.Definition, t.ID,
	)
	if err != nil {
		return fmt.Errorf("UpdateGlossaryTerm: %w", err)
	}
	return nil
}

func (r *DBRepository) DeleteGlossaryTerm(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, nil, `DELETE FROM glossary_terms WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("DeleteGlossaryTerm: %w", err)
	}
	return nil
}

func (r *DBRepository) ListStudentProgress(ctx context.Context, userID int64) ([]dbEntities.StudentProgress, error) {
	var progress []dbEntities.StudentProgress
	err := r.db.Select(ctx, nil, &progress,
		`SELECT sp.id, sp.user_id, COALESCE(NULLIF(u.display_name, ''), u.login) AS user_display_name,
		        sp.discipline, sp.level, sp.passed_checks, sp.instructor_note, sp.updated_by, sp.updated_at
		 FROM student_progress sp
		 JOIN users u ON u.id=sp.user_id
		 WHERE ($1 = 0 OR sp.user_id=$1)
		 ORDER BY u.display_name ASC, sp.discipline ASC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("ListStudentProgress: %w", err)
	}
	return progress, nil
}

func (r *DBRepository) UpsertStudentProgress(ctx context.Context, p *dbEntities.StudentProgress) error {
	err := r.db.QueryRow(ctx, nil,
		[]any{&p.ID},
		`INSERT INTO student_progress
		   (user_id, discipline, level, passed_checks, instructor_note, updated_by, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,NOW())
		 ON CONFLICT (user_id, discipline) DO UPDATE SET
		   level=EXCLUDED.level, passed_checks=EXCLUDED.passed_checks,
		   instructor_note=EXCLUDED.instructor_note, updated_by=EXCLUDED.updated_by,
		   updated_at=NOW()
		 RETURNING id`,
		p.UserID, p.Discipline, p.Level, p.PassedChecks, p.InstructorNote, p.UpdatedBy,
	)
	if err != nil {
		return fmt.Errorf("UpsertStudentProgress: %w", err)
	}
	return nil
}
