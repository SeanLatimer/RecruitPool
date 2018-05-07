package modules

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc"
	"github.com/sirupsen/logrus"
)

var (
	RecruitPool = &recruitPool{
		subOnly: true,
	}
)

type recruitPool struct {
	pool    []string
	subOnly bool
	log     *logrus.Logger
	open    bool
}

func (r *recruitPool) PoolAdd(user twitch.User, sub bool) {
	if user.UserType != "mod" {
		if r.subOnly && !sub {
			return
		}
	}
	r.pool = append(r.pool, user.Username)
}

func (r *recruitPool) PoolOpen() {
	r.open = true
}

func (r *recruitPool) PoolClose() {
	r.open = false
}

func (r *recruitPool) PoolIsOpen() bool {
	return r.open
}

func (r *recruitPool) Draw() string {
	if len(r.pool) == 0 {
		return ""
	}

	rand.Seed(time.Now().Unix())
	r.pool = dedup(r.pool)
	winner := r.pool[rand.Intn(len(r.pool))]
	r.pool = remove(r.pool, winner)

	return winner
}

func (r *recruitPool) Count() int {
	r.log.Info(len(r.pool))
	return len(r.pool)
}

func (r *recruitPool) PoolReset() {
	r.pool = make([]string, 0)
}

func (r *recruitPool) PoolSetSubOnly(state bool) {
	r.subOnly = state
}

func (r *recruitPool) PoolIsSubOnly() bool {
	return r.subOnly
}

func (r *recruitPool) HandleCommands(args string, user twitch.User, messageTags map[string]string, broadcaster bool) string {
	if broadcaster {
		user.UserType = "mod"
		switch {
		case strings.HasPrefix(strings.ToLower(args), "open"):
			if !r.open {
				r.open = true
				if len(r.pool) > 0 {
					return `Pool Open with existing recruits! You can reset with !recruit reset. Type !recruit to join`
				}
				return `Pool Open! Type !recruit to join`
			}
		case strings.HasPrefix(strings.ToLower(args), "close"):
			if r.open {
				r.open = false
				return "Pool is now closed!"
			}
		case strings.HasPrefix(strings.ToLower(args), "draw"):
			winner := r.Draw()
			if winner == "" {
				return "Oh no, looks like we ran out of recruits! Type !recruit to join"
			}
			return "Congratulations to our new recruit " + winner + "!"
		case strings.HasPrefix(strings.ToLower(args), "subonly"):
			if !r.PoolIsSubOnly() {
				r.PoolSetSubOnly(!r.PoolIsSubOnly())
				return "Pool is now sub only!"
			}
			r.PoolSetSubOnly(!r.PoolIsSubOnly())
			return "Pool is no longer sub only!"

		case strings.HasPrefix(strings.ToLower(args), "reset"):
			r.PoolReset()
			return "Pool reset!"
		case strings.HasPrefix(strings.ToLower(args), "count"):
			return "We have " + strconv.Itoa(r.Count()) + " recruits remaining!"
		}
	}

	sub, _ := strconv.ParseBool(messageTags["subscriber"])
	r.PoolAdd(user, sub)

	return ""
}

func (r *recruitPool) GetBaseCommand() string {
	return "recruit"
}

func (r *recruitPool) SetLogger(logger *logrus.Logger) {
	r.log = logger
}

func dedup(slice []string) []string {
	var returnSlice []string
	for _, value := range slice {
		if !contains(returnSlice, value) {
			returnSlice = append(returnSlice, value)
		}
	}
	return returnSlice
}

func contains(slice []string, s string) bool {
	for _, vv := range slice {
		if vv == s {
			return true
		}
	}
	return false
}

func remove(slice []string, remove string) []string {
	var newPool []string
	for _, value := range slice {
		if value != remove {
			newPool = append(newPool, value)
		}
	}

	return newPool
}
