# String -> Query Parser for Tormenta

Turn (URL) strings like this:

```
query=limit:1,offset:1,reverse:true,index:indexstring,start:1,end:10&query=index:indexstring,start:1,end:100
```

into Tormenta queries like:

```golang
db.And(
    &results,
    db.Find(&results).Limit(1).Offset(1).Reverse().Range("indexstring", 1, 10)
    db.Find(&results).Range("indexstring", 1, 100)
)
```

## Notes

- `AND` combination by default, use `or=true` for OR
- Cannot deal with complex nesting of AND/OR yet