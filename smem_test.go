package smem

import (
	"sync"
	"testing"
	"time"
)

func Test_smem(t *testing.T) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		store := NewStore()

		expire := time.Now().Add(1 * time.Hour)
		newSession, err := store.NewSession("test", expire)
		if err != nil {
			t.Fatal(err)
		}

		session, ok := store.ReadSession(newSession.Token())
		if !ok {
			t.Fatal("Expect to find session", newSession.Token())
		}

		if session.Token() != newSession.Token() ||
			session.Expires() != expire ||
			session.UserID() != "test" {
			t.Fatal("Wrong data in returned session", session)
		}

		_, ok = store.ReadSession("none")
		if ok {
			t.Fatal("Expect to find no session")
		}

		err = store.RemoveSession(newSession.Token())
		if err != nil {
			t.Fatal(err)
		}

		out := time.Now().Add(-1 * time.Hour)
		newSession, err = store.NewSession("test2", out)

		session, ok = store.ReadSession(newSession.Token())
		if ok {
			t.Fatal("Expect to find no session due tu expired")
		}

		out = time.Now().Add(-1 * time.Hour)
		_, err = store.NewSession("test3", out)

		n, err := store.RemoveExpiredSessions()
		if err != nil {
			t.Fatal(err)
		}

		if n != 1 {
			t.Fatal("Expect to remove one")
		}

	}()

	wg.Wait()
}
