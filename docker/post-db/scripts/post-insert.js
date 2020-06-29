function get_results(result) {
    print(tojson(result));
}

function insert_post(object) {
    print(db.posts.insert(object));
}

insert_post({
    "_id": ObjectId("5eeb690cee11cb025b774bb5"),
    "kamigo": "運動会",
    "nakashichi": "抜くなその子は",
    "shimogo": "課長の子",
    "user_id": ObjectId("57a98d98e4b00679b4a830af"),
    "signup_at": ISODate("2020-06-27 18:00:00+09:00")
});

print("_______POST DATA_______");
db.posts.find().forEach(get_results);
print("______END POST DATA_____");
