package main

import "testing"

func TestRemoveExtension(t *testing.T) {

	assertEquals := func(t *testing.T, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("expected '%q' but got '%q'", got, want)
		}
	}
	t.Run("jpeg", func(t *testing.T) {
		got := removeExtension("test.jpeg")
		want := "test"
		assertEquals(t, got, want)
	})
	t.Run("png", func(t *testing.T) {
		got := removeExtension("test.jpeg")
		want := "test"
		assertEquals(t, got, want)
	})
	t.Run("some_imaginary_format", func(t *testing.T) {
		got := removeExtension("test.jpeg")
		want := "test"
		assertEquals(t, got, want)
	})

}
