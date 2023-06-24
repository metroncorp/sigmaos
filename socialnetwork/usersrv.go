package socialnetwork

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"sigmaos/cacheclnt"
	"sigmaos/mongoclnt"
	dbg "sigmaos/debug"
	"sigmaos/perf"
	"sigmaos/fs"
	"sigmaos/protdevsrv"
	sp "sigmaos/sigmap"
	"gopkg.in/mgo.v2/bson"
	"sigmaos/socialnetwork/proto"
	"sync"
)

// YH:
// User service for social network
// for now we use sql instead of MongoDB

const (
	USER_QUERY_OK = "OK"
	USER_CACHE_PREFIX = "user_"
)

type UserSrv struct {
	mu     sync.Mutex
	mongoc *mongoclnt.MongoClnt
	cachec *cacheclnt.CacheClnt
	sid    int32 // sid is a random number between 0 and 2^30
	ucount int32 //This server may overflow with over 2^31 users
}

func RunUserSrv(public bool, jobname string) error {
	dbg.DPrintf(dbg.SOCIAL_NETWORK_USER, "Creating user service\n")
	usrv := &UserSrv{}
	usrv.sid = rand.Int31n(536870912) // 2^29
	pds, err := protdevsrv.MakeProtDevSrvPublic(sp.SOCIAL_NETWORK_USER, usrv, public)
	if err != nil {
		return err
	}
	mongoc, err := mongoclnt.MkMongoClnt(pds.MemFs.SigmaClnt().FsLib)
	if err != nil {
		return err
	}
	mongoc.EnsureIndex(SN_DB, USER_COL, []string{"userid"})
	usrv.mongoc = mongoc
	fsls := MakeFsLibs(sp.SOCIAL_NETWORK_USER)
	cachec, err := cacheclnt.MkCacheClnt(fsls, jobname)
	if err != nil {
		return err
	}
	usrv.cachec = cachec
	dbg.DPrintf(dbg.SOCIAL_NETWORK_USER, "Starting user service %v\n", usrv.sid)
	perf, err := perf.MakePerf(perf.SOCIAL_NETWORK_USER)
	if err != nil {
		dbg.DFatalf("MakePerf err %v\n", err)
	}
	defer perf.Done()
	return pds.RunServer()
}

func (usrv *UserSrv) CheckUser(ctx fs.CtxI, req proto.CheckUserRequest, res *proto.CheckUserResponse) error {
	dbg.DPrintf(dbg.SOCIAL_NETWORK_USER, "Checking user at %v: %v\n", usrv.sid, req.Usernames)
	userids := make([]int64, len(req.Usernames))
	res.Ok = "No"
	missing := false
	for idx, username := range req.Usernames {
		user, err := usrv.getUserbyUname(username)
		if err != nil {
			return err
		}
		if user == nil {
			userids[idx] = int64(-1)
			missing = true
		} else {
			userids[idx] = user.Userid
		}
	}
	res.Userids = userids
	if !missing {
		res.Ok = USER_QUERY_OK
	}
	return nil
}

func (usrv *UserSrv) RegisterUser(ctx fs.CtxI, req proto.RegisterUserRequest, res *proto.UserResponse) error {
	dbg.DPrintf(dbg.SOCIAL_NETWORK_USER, "Register user at %v: %v\n", usrv.sid, req)
	res.Ok = "No"
	user, err := usrv.getUserbyUname(req.Username)
	if err != nil {
		return err
	}
	if user != nil {
		res.Ok = fmt.Sprintf("Username %v already exist", req.Username)
		return nil
	}
	pswd_hashed := fmt.Sprintf("%x", sha256.Sum256([]byte(req.Password)))
	userid := usrv.getNextUserId()
	newUser := User{
		Userid: userid,
		Username: req.Username,
		Lastname: req.Lastname,
		Firstname: req.Firstname,
		Password: pswd_hashed}
	if err := usrv.mongoc.Insert(SN_DB, USER_COL, newUser); err != nil {
		dbg.DFatalf("Mongo Error: %v", err)
		return err
	}
	res.Ok = USER_QUERY_OK
	res.Userid = userid
	return nil
}

func (usrv *UserSrv) incCountSafe() int32 {
	usrv.mu.Lock()
	defer usrv.mu.Unlock()
	usrv.ucount++
	return usrv.ucount
}

func (usrv *UserSrv) getNextUserId() int64 {
	return int64(usrv.sid)*1e10 + int64(usrv.incCountSafe())
}

func (usrv *UserSrv) Login(ctx fs.CtxI, req proto.LoginRequest, res *proto.UserResponse) error {
	dbg.DPrintf(dbg.SOCIAL_NETWORK_USER, "User login with %v: %v\n", usrv.sid, req)
	res.Ok = "Login Failure."
	user, err := usrv.getUserbyUname(req.Username)
	if err != nil {
		return err
	}
	if user != nil && fmt.Sprintf("%x", sha256.Sum256([]byte(req.Password))) == user.Password {
		res.Ok = USER_QUERY_OK
		res.Userid = user.Userid
	}
	return nil
}

func (usrv *UserSrv) checkUserExist(username string) (bool, error) {
	user, err := usrv.getUserbyUname(username)
	if err != nil {
		return false, err
	}
	return user != nil, nil
}

func (usrv *UserSrv) getUserbyUname(username string) (*User, error) {
	key := USER_CACHE_PREFIX + username
	user := &User{}
	cacheItem := &proto.CacheItem{}
	if err := usrv.cachec.Get(key, cacheItem); err != nil {
		if !usrv.cachec.IsMiss(err) {
			return nil, err
		}
		dbg.DPrintf(dbg.SOCIAL_NETWORK_USER, "User %v cache miss\n", key)
		found, err := usrv.mongoc.FindOne(SN_DB, USER_COL, bson.M{"username": username}, user)
		if err != nil {
			return nil, err
		} 
		if !found {
			return nil, nil
		}
		encoded, _ := bson.Marshal(user)
		usrv.cachec.Put(key, &proto.CacheItem{Key: key, Val: encoded})
		dbg.DPrintf(dbg.SOCIAL_NETWORK_USER, "Found user %v in DB: %v\n", username, user)
	} else {
		bson.Unmarshal(cacheItem.Val, user)
		dbg.DPrintf(dbg.SOCIAL_NETWORK_USER, "Found user %v in cache!\n", username)
	}
	return user, nil
}

type User struct {
	Userid    int64  `bson:userid`
	Firstname string `bson:firstname`
	Lastname  string `bson:lastname`
	Username  string `bson:username`
	Password  string `bson:password`
}
