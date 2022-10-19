# ent gqlgen todos

## 作業手順

### base

1. schemaファイル作成する
    ```shell
    go run -mod=mod entgo.io/ent/cmd/ent init <Model>
    ```
    - `./ent/schema/<model>.go` が生成されます。
2. `Fields()`と`Edges()`を定義する
    - 対象ファイル: `./ent/schema/<model>.go`
    - 参照:
      - [Fields](https://entgo.io/ja/docs/schema-fields)
      - [Edges](https://entgo.io/ja/docs/schema-edges)
3. アセットを生成する
    ```shell
    go generate .
    ```

### Query

1. Query リゾルバーを実装する
   - `./graph/ent.resolvers.go`
 
### Mutation

1. `<model>.graphql`を作成する
    ```shell
    touch ./graph/<model>.graphql
    ```
2. Mutationを定義する
3. コードを生成する
    ```shell
    go generate .
    ```
4. Mutation リゾルバーを実装する
    - `./graph/<model>.resolvers.go`

### Pagination

- [Pagination](https://entgo.io/ja/docs/tutorial-todo-gql-paginate)

### Mutation Inputs

- [Mutation Inputs](https://entgo.io/ja/docs/tutorial-todo-gql-mutation-input)

### Filter Inputs

- [Filter Inputs](https://entgo.io/ja/docs/tutorial-todo-gql-filter-input)

## スキーマの詳細を確認する

```shell
go run -mod=mod entgo.io/ent/cmd/ent describe ./ent/schema
```
