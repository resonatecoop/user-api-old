package model

import (
  "time"

  "github.com/satori/go.uuid"
  "github.com/go-pg/pg"
)

// AuthUser represents data stored in session/context for a user
type AuthUser struct {
	Id       uuid.UUID
	TenantId int32
	Username string
	Email    string
	Role     AccessRole
}

type User struct {
  Id uuid.UUID `sql:"type:uuid,default:uuid_generate_v4()"`
  CreatedAt time.Time `sql:"default:now()"`
  UpdatedAt time.Time
  Username string `sql:",notnull,unique"`
  FullName string `sql:",notnull"`
  FirstName string
  LastName string
  Email string `sql:",unique,notnull"`
  Member bool `sql:",notnull"`
  NewsletterNotification bool

  FavoriteTracks []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  FollowedGroups []uuid.UUID `sql:",type:uuid[]" pg:",array"`
  Playlists []uuid.UUID `sql:",type:uuid[]" pg:",array"`

  OwnerOfGroups []UserGroup `pg:"fk:owner_id"`
  // Plays []Track `pg:"many2many:plays"` Payment API

  TenantId int32
  RoleId int32
  LastLogin          *time.Time
  LastPasswordChange *time.Time
  Password string
  Token string
}

// UpdateLoginDetails updates login related fields
func (u *User) UpdateLoginDetails(token string) {
	u.Token = token
	t := time.Now()
	u.LastLogin = &t
}

func (u *User) Delete(tx *pg.Tx) (error, string) {
  pgerr := tx.Model(u).
    Column("user.favorite_tracks", "user.followed_groups", "user.playlists", "OwnerOfGroups").
    WherePK().
    Select()
  if pgerr != nil {
    return pgerr, "user"
  }

  if len(u.FavoriteTracks) > 0 {
    _, pgerr = tx.Exec(`
      UPDATE tracks
      SET favorite_of_users = array_remove(favorite_of_users, ?)
      WHERE id IN (?)
    `, u.Id, pg.In(u.FavoriteTracks))
    if pgerr != nil {
      return pgerr, "track"
    }
  }

  if len(u.FollowedGroups) > 0 {
    _, pgerr = tx.Exec(`
      UPDATE user_groups
      SET followers = array_remove(followers, ?)
      WHERE id IN (?)
    `, u.Id, pg.In(u.FollowedGroups))
    if pgerr != nil {
      return pgerr, "user_group"
    }
  }

  if len(u.OwnerOfGroups) > 0 {
    for _, group := range u.OwnerOfGroups {
      if pgerr, table := group.Delete(tx); pgerr != nil {
        return pgerr, table
      }
    }
  }

  if len(u.Playlists) > 0 {
    for _, playlistId := range u.Playlists {
      playlist := &TrackGroup{Id: playlistId}
      if pgerr, table := playlist.Delete(tx); pgerr != nil {
        return pgerr, table
      }
    }
  }

  pgerr = tx.Delete(u)
  if pgerr != nil {
    return pgerr, "user"
  }

  return nil, ""
}

func (u *User) FollowGroup(db *pg.DB, userGroupId uuid.UUID) (error, string) {
  var table string
  tx, err := db.Begin()
  if err != nil {
    return err, table
  }
  defer tx.Rollback()

  // Add userGroupId to user FollowedGroups
  userGroupIdArr := []uuid.UUID{userGroupId}
  _, pgerr := tx.ExecOne(`
    UPDATE users
    SET followed_groups = (select array_agg(distinct e) from unnest(followed_groups || ?) e)
    WHERE id = ?
  `, pg.Array(userGroupIdArr), u.Id)
  if pgerr != nil {
    table = "user"
    return pgerr, table
  }

  // Add userId to userGroup Followers
  userIdArr := []uuid.UUID{u.Id}
  _, pgerr = tx.ExecOne(`
    UPDATE user_groups
    SET followers = (select array_agg(distinct e) from unnest(followers || ?) e)
    WHERE id = ?
  `, pg.Array(userIdArr), userGroupId)
  if pgerr != nil {
    table = "user_group"
    return pgerr, table
  }
  return tx.Commit(), table
}

func (u *User) UnfollowGroup(db *pg.DB, userGroupId uuid.UUID) (error, string) {
  var table string
  tx, err := db.Begin()
  if err != nil {
    return err, table
  }
  // Rollback tx on error.
  defer tx.Rollback()

  // Remove userGroupId from user FollowedGroups
  _, pgerr := tx.ExecOne(`
    UPDATE users
    SET followed_groups = array_remove(followed_groups, ?)
    WHERE id = ?
  `, userGroupId, u.Id)
  if pgerr != nil {
    table = "user"
    return pgerr, table
  }

  // Remove userId from track FavoriteOfUsers
  _, pgerr = tx.ExecOne(`
    UPDATE user_groups
    SET followers = array_remove(followers, ?)
    WHERE id = ?
  `, u.Id, userGroupId)
  if pgerr != nil {
    table = "user_group"
    return pgerr, table
  }
  return tx.Commit(), table
}

func (u *User) RemoveFavoriteTrack(db *pg.DB, trackId uuid.UUID) (error, string) {
  var table string
  tx, err := db.Begin()
  if err != nil {
    return err, table
  }
  // Rollback tx on error.
  defer tx.Rollback()

  // Remove trackId from user FavoriteTracks
  _, pgerr := tx.ExecOne(`
    UPDATE users
    SET favorite_tracks = array_remove(favorite_tracks, ?)
    WHERE id = ?
  `, trackId, u.Id)
  if pgerr != nil {
    table = "user"
    return pgerr, table
  }

  // Remove userId from track FavoriteOfUsers
  _, pgerr = tx.ExecOne(`
    UPDATE tracks
    SET favorite_of_users = array_remove(favorite_of_users, ?)
    WHERE id = ?
  `, u.Id, trackId)
  if pgerr != nil {
    table = "track"
    return pgerr, table
  }
  return tx.Commit(), table
}

func (u *User) AddFavoriteTrack(db *pg.DB, trackId uuid.UUID) (error, string) {
  var table string
  tx, err := db.Begin()
  if err != nil {
    return err, table
  }
  // Rollback tx on error.
  defer tx.Rollback()

  // Add trackId to user FavoriteTracks
  trackIdArr := []uuid.UUID{trackId}
  _, pgerr := tx.ExecOne(`
    UPDATE users
    SET favorite_tracks = (select array_agg(distinct e) from unnest(favorite_tracks || ?) e)
    WHERE id = ?
  `, pg.Array(trackIdArr), u.Id)
  // WHERE NOT favorite_tracks @> ?
  if pgerr != nil {
    table = "user"
    return pgerr, table
  }

  // Add userId to track FavoriteOfUsers
  userIdArr := []uuid.UUID{u.Id}
  _, pgerr = tx.ExecOne(`
    UPDATE tracks
    SET favorite_of_users = (select array_agg(distinct e) from unnest(favorite_of_users || ?) e)
    WHERE id = ?
  `, pg.Array(userIdArr), trackId)
  if pgerr != nil {
    table = "track"
    return pgerr, table
  }
  return tx.Commit(), table
}
