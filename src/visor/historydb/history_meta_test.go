package historydb

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/samoslab/galtcoin/src/visor/dbutil"
)

func TestHistoryMetaGetSetParsedHeight(t *testing.T) {
	db, td := prepareDB(t)
	defer td()

	hm := &historyMeta{}

	err := db.View("", func(tx *dbutil.Tx) error {
		height, ok, err := hm.ParsedBlockSeq(tx)
		require.NoError(t, err)
		require.False(t, ok)
		require.Equal(t, uint64(0), height)
		return err
	})
	require.NoError(t, err)

	err = db.Update("", func(tx *dbutil.Tx) error {
		err := hm.SetParsedBlockSeq(tx, 10)
		require.NoError(t, err)
		return err
	})
	require.NoError(t, err)

	err = db.View("", func(tx *dbutil.Tx) error {
		height, ok, err := hm.ParsedBlockSeq(tx)
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, uint64(10), height)
		return err
	})
	require.NoError(t, err)

	err = db.Update("", func(tx *dbutil.Tx) error {
		err := hm.SetParsedBlockSeq(tx, 0)
		require.NoError(t, err)
		return err
	})
	require.NoError(t, err)

	err = db.View("", func(tx *dbutil.Tx) error {
		height, ok, err := hm.ParsedBlockSeq(tx)
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, uint64(0), height)
		return err
	})
	require.NoError(t, err)

}
