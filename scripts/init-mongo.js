db.getSiblingDB("admin").auth("admin", "pass");
db.createUser({
  user: "user",
  pwd: "pass",
  roles: ["readWrite"],
});
