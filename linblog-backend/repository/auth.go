package repository

import "time"

func (r *Repository) CountDownKey(key string, value string, expiration time.Time) error {
	return r.rc.Set(key, value, time.Until(expiration)).Err()
}
