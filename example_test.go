package todo

import (
	"context"
	"fmt"
	"github.com/shiroemons/ent-gqlgen-todos/ent/todo"
	"log"

	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3"

	"github.com/shiroemons/ent-gqlgen-todos/ent"
)

func Example_Todo() {
	// インメモリーのSQLiteデータベースを持つent.Clientを作成します。
	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	ctx := context.Background()
	// 自動マイグレーションツールを実行して、すべてのスキーマリソースを作成します。
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	// 出力します。
	task1, err := client.Todo.Create().SetText("Add GraphQL Example").Save(ctx)
	if err != nil {
		log.Fatalf("failed creating a todo: %v", err)
	}
	task2, err := client.Todo.Create().SetText("Add Tracing Example").Save(ctx)
	if err != nil {
		log.Fatalf("failed creating a todo: %v", err)
	}
	if err := task2.Update().SetParent(task1).Exec(ctx); err != nil {
		log.Fatalf("failed connecting todo2 to its parent: %v", err)
	}
	// 子TODOを通じて親TODOを取得し、
	// クエリが正確に1つのTODOを返すことを期待します。
	parent, err := client.Todo.Query(). // すべてのtodoアイテムを取得する
		Where(todo.HasParent()). // 親todoアイテムを持つtodoアイテムのみにフィルタリング
		QueryParent(). // 親todoアイテムについて走査を続ける
		Only(ctx) // 1つのtodoアイテムのみ取得する
	if err != nil {
		log.Fatalf("failed querying todos: %v", err)
	}
	fmt.Printf("%d: %q\n", parent.ID, parent.Text)
	// Output:
	// 1: "Add GraphQL Example"
}
