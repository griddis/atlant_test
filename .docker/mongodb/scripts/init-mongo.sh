mongo -- "$MONGO_DB" <<EOF
var rootuser = '$MONGO_INITDB_ROOT_USERNAME';
var rootpasswd = '$MONGO_INITDB_ROOT_PASSWORD';
var admin = db.getSiblingDB('admin');
admin.auth(rootuser, rootpasswd);
db.createUser({user: rootuser, pwd: rootpasswd, roles: ["readWrite"]});
EOF