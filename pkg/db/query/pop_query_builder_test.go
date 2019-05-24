package query

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type PopQueryBuilderSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *PopQueryBuilderSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestUserSuite(t *testing.T) {

	hs := &PopQueryBuilderSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, hs)
}

type testQueryFilter struct {
	column     string
	comparator string
	value      string
}

func (f testQueryFilter) Column() string {
	return f.column
}

func (f testQueryFilter) Comparator() string {
	return f.comparator
}

func (f testQueryFilter) Value() string {
	return f.value
}

func (suite *PopQueryBuilderSuite) TestFetchOne() {
	user := testdatagen.MakeDefaultOfficeUser(suite.DB())
	builder := NewPopQueryBuilder(suite.DB())
	var actualUser models.OfficeUser

	suite.T().Run("fetches one with filter", func(t *testing.T) {
		// create extra record to make sure we filter
		user2 := testdatagen.MakeDefaultOfficeUser(suite.DB())
		filter := testQueryFilter{"id", equals, user.ID.String()}

		err := builder.FetchOne(&actualUser, filter)

		suite.NoError(err)
		suite.Equal(user.ID, actualUser.ID)

		// do the reverse to make sure we don't get the same record every time
		filter = testQueryFilter{"id", equals, user2.ID.String()}

		err = builder.FetchOne(&actualUser, filter)

		suite.NoError(err)
		suite.Equal(user2.ID, actualUser.ID)
	})

	suite.T().Run("returns error on invalid column", func(t *testing.T) {
		filter := testQueryFilter{"fake_column", equals, user.ID.String()}
		var actualUser models.OfficeUser

		err := builder.FetchOne(&actualUser, filter)

		suite.Error(err)
		suite.Equal("[fake_column] is not valid input", err.Error())
		suite.Zero(actualUser)
	})

	suite.T().Run("fails when not pointer", func(t *testing.T) {
		var actualUser models.OfficeUser

		err := builder.FetchOne(actualUser)

		suite.Error(err)
		suite.Equal("Model should be pointer to struct", err.Error())
		suite.Zero(actualUser)
	})

	suite.T().Run("fails when not pointer to struct", func(t *testing.T) {
		var i int

		err := builder.FetchOne(&i)

		suite.Error(err)
		suite.Equal("Model should be pointer to struct", err.Error())
	})

}

func (suite *PopQueryBuilderSuite) TestFetchMany() {
	// this should be stubbed out with a model that is agnostic to our code
	// similar to how the pop repo tests might work
	user := testdatagen.MakeDefaultOfficeUser(suite.DB())
	builder := NewPopQueryBuilder(suite.DB())
	var actualUsers models.OfficeUsers

	suite.T().Run("fetches many with filter", func(t *testing.T) {
		user2 := testdatagen.MakeDefaultOfficeUser(suite.DB())
		filter := testQueryFilter{"id", equals, user2.ID.String()}

		err := builder.FetchMany(&actualUsers, filter)

		suite.NoError(err)
		suite.Len(actualUsers, 1)
		suite.Equal(user2.ID, actualUsers[0].ID)

		// do the reverse to make sure we don't get the same record every time
		filter = testQueryFilter{"id", equals, user.ID.String()}
		var actualUsers models.OfficeUsers

		err = builder.FetchMany(&actualUsers, filter)

		suite.NoError(err)
		suite.Len(actualUsers, 1)
		suite.Equal(user.ID, actualUsers[0].ID)
	})

	suite.T().Run("fails with invalid column", func(t *testing.T) {
		var actualUsers models.OfficeUsers
		filter := testQueryFilter{"fake_column", equals, user.ID.String()}

		err := builder.FetchMany(&actualUsers, filter)

		suite.Error(err)
		suite.Equal("[fake_column] is not valid input", err.Error())
		suite.Empty(actualUsers)
	})

	suite.T().Run("fails when not pointer", func(t *testing.T) {
		var actualUsers models.OfficeUsers

		err := builder.FetchMany(actualUsers)

		suite.Error(err)
		suite.Equal("Model should be pointer to slice of structs", err.Error())
		suite.Empty(actualUsers)
	})

	suite.T().Run("fails when not pointer to slice", func(t *testing.T) {
		var actualUser models.OfficeUser

		err := builder.FetchMany(&actualUser)

		suite.Error(err)
		suite.Equal("Model should be pointer to slice of structs", err.Error())
		suite.Empty(actualUser)
	})

	suite.T().Run("fails when not pointer to slice of structs", func(t *testing.T) {
		var intSlice []int

		err := builder.FetchMany(&intSlice)

		suite.Error(err)
		suite.Equal("Model should be pointer to slice of structs", err.Error())
	})
}