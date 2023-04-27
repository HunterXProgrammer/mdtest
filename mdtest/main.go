// Copyright (c) 2021 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os/exec"
	"io"
	"io/ioutil"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"google.golang.org/protobuf/proto"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/appstate"
	waBinary "go.mau.fi/whatsmeow/binary"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)
var server *http.Server
var is_mode string
var tasker_type string
var is_send bool
var is_connected bool = true
type test_struct struct {
    Data1 string `json:"data1"`
    Data2 string `json:"data2"`
    Data3 string `json:"data3"`
    Data4 string `json:"data4"`
    Data5 string `json:"data5"`
    Data6 string `json:"data6"`
    Data7 string `json:"data7"`
    Data8 string `json:"data8"`
    Data9 string `json:"data9"`
    Data10 string `json:"data10"`
    Data11 string `json:"data11"`
    Data12 string `json:"data12"`
    Data13 string `json:"data13"`
    Data14 string `json:"data14"`
    Data15 string `json:"data15"`
    Data16 string `json:"data16"`
    Data17 string `json:"data17"`
    Data18 string `json:"data18"`
    Data19 string `json:"data19"`
    Data20 string `json:"data20"`
    Data21 string `json:"data21"`
    Data22 string `json:"data22"`
    Data23 string `json:"data23"`
    Data24 string `json:"data24"`
    Data25 string `json:"data25"`
    Data26 string `json:"data26"`
    Data27 string `json:"data27"`
    Data28 string `json:"data28"`
    Data29 string `json:"data29"`
    Data30 string `json:"data30"`
    Data31 string `json:"data31"`
    Data32 string `json:"data32"`
    Data33 string `json:"data33"`
    Data34 string `json:"data34"`
    Data35 string `json:"data35"`
    Data36 string `json:"data36"`
    Data37 string `json:"data37"`
    Data38 string `json:"data38"`
    Data39 string `json:"data39"`
    Data40 string `json:"data40"`
}

func mdtest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		log.Errorf("Invalid request path, 404 not found.")
		return
	}

	switch r.Method {
	case "GET":
		if is_connected {
			if is_mode == "both" {
		    	log.Infof("get request received, mdtest is running in both mode")
		    	io.WriteString(w, "mdtest is running in both mode")
	    	} else if is_mode == "receive" {
	    		log.Infof("get request received, mdtest is running in receive mode")
		    	io.WriteString(w, "mdtest is running in receive mode")
    		} else if is_mode == "send" {
    			log.Infof("get request received, mdtest is running in send mode")
		    	io.WriteString(w, "mdtest is running in send mode")
			}
	    } else {
	    	log.Infof("get request received, mdtest is waiting for reconnection")
		    io.WriteString(w, "Bad network, mdtest is waiting for reconnection")
    	}
	case "POST":
  	  dec := json.NewDecoder(r.Body)
		    for {
			var t test_struct
    	    if err := dec.Decode(&t); err == io.EOF {
   	         break
    	    } else if err != nil {
    	        log.Errorf("%s", err)
    	    }
   	     
   	     if t.Data1 == "stop" {
   	     	io.WriteString(w, "exiting")
				go func() {
					time.Sleep(1 * time.Second)
					kill_server()
					log.Infof("Exit command Received, exiting...")
					cli.Disconnect()
					os.Exit(0)
				}()
				return
			} else if t.Data1 == "disconnect" {
				io.WriteString(w, "disconnected")
				go func() {
					if is_connected {
						is_connected = false
						log.Infof("Bad network, mdtest is waiting for reconnection")
					}
				}()
				return
			}  else if t.Data1 == "connect" {
				io.WriteString(w, "connected")
				go func() {
					if !is_connected {
						is_connected = true
						log.Infof("Network restored, mdtest is running")
					}
				}()
				return
			}
   	     io.WriteString(w, "command received")
   	     if t.Data1 != "" && t.Data2 == "" {
				log.Infof("no data2")
				args := []string{t.Data1, "null"}
				if is_mode == "both" ||  is_mode == "send" {
					go handleCmd(strings.ToLower(args[0]), args[1:])
				}
				return
			} else if t.Data1 != "" && t.Data2 != "" {
				args := []string{t.Data1, t.Data2, t.Data3, t.Data4, t.Data5, t.Data6, t.Data7, t.Data8, t.Data9, t.Data10, t.Data11, t.Data12, t.Data13, t.Data14, t.Data15, t.Data16, t.Data17, t.Data18, t.Data19, t.Data20, t.Data21, t.Data22, t.Data23, t.Data24, t.Data25, t.Data26, t.Data27, t.Data28, t.Data29, t.Data30, t.Data31, t.Data32, t.Data33, t.Data34, t.Data35, t.Data36, t.Data37, t.Data38, t.Data39, t.Data40}
				if is_mode == "both" ||  is_mode == "send" {
					go handleCmd(strings.ToLower(args[0]), args[1:])
				}
				return
			}
		}
	    
	default:
		log.Errorf("%s, only GET and POST methods are supported.", w)
	}
}

func kill_server() {
	//if is_mode == "receive" {
		//return
	//}
	resp, err := http.Get("http://localhost:7777")
	if err == nil {
	    body, err := ioutil.ReadAll(resp.Body)
	    resp.Body.Close()
    	if err == nil {
    	sb := string(body)
    
        	if strings.Contains(sb, "mdtest is running") || strings.Contains(sb, "Bad network") {
            	server.Close()
        	}
    	}
	}
}

func MdtestStart() {
	//if is_mode == "receive" {
		//return
	//}
	
	resp, err := http.Get("http://localhost:7777")
	if err == nil {
    	body, err := ioutil.ReadAll(resp.Body)
    	resp.Body.Close()
    	if err == nil {
    	sb := string(body)
    
        	if !strings.Contains(sb, "mdtest is running") && !strings.Contains(sb, "Bad network") {
            	http.HandleFunc("/", mdtest)
				log.Infof("mdtest started")
				server = &http.Server{
					Addr:    ":7777",
        			Handler: http.DefaultServeMux,
    			}
    			server.ListenAndServe()
    			log.Infof("mdtest stopped")
        	}
    	}
	} else {
    	http.HandleFunc("/", mdtest)
		log.Infof("mdtest started")
		server = &http.Server{
			Addr:    ":7777",
			Handler: http.DefaultServeMux,
		}
		server.ListenAndServe()
		log.Infof("mdtest stopped")
	}
}

var cli *whatsmeow.Client
var log waLog.Logger

var logLevel = "INFO"
var debugLogs = flag.Bool("debug", false, "Enable debug logs?")
var dbDialect = flag.String("db-dialect", "sqlite3", "Database dialect (sqlite3 or postgres)")
var dbAddress = flag.String("db-address", "file:mdtest.txt?_foreign_keys=on", "Database address")
var requestFullSync = flag.Bool("request-full-sync", false, "Request full (1 year) history sync when logging in?")
var pairRejectChan = make(chan bool, 1)

func main() {
	waBinary.IndentXML = true
	flag.Parse()

	if *debugLogs {
		logLevel = "DEBUG"
	}
	if *requestFullSync {
		store.DeviceProps.RequireFullSync = proto.Bool(true)
	}
	log = waLog.Stdout("Main", logLevel, true)

	dbLog := waLog.Stdout("Database", logLevel, true)
	storeContainer, err := sqlstore.New(*dbDialect, *dbAddress, dbLog)
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
		return
	}
	device, err := storeContainer.GetFirstDevice()
	if err != nil {
		log.Errorf("Failed to get device: %v", err)
		return
	}

	cli = whatsmeow.NewClient(device, waLog.Stdout("Client", logLevel, true))
	var isWaitingForPair atomic.Bool
	cli.PrePairCallback = func(jid types.JID, platform, businessName string) bool {
		isWaitingForPair.Store(true)
		defer isWaitingForPair.Store(false)
		log.Infof("Pairing %s (platform: %q, business name: %q). Type r within 3 seconds to reject pair", jid, platform, businessName)
		select {
		case reject := <-pairRejectChan:
			if reject {
				log.Infof("Rejecting pair")
				return false
			}
		case <-time.After(3 * time.Second):
		}
		log.Infof("Accepting pair")
		return true
	}

	ch, err := cli.GetQRChannel(context.Background())
	if err != nil {
		// This error means that we're already logged in, so ignore it.
		if !errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
			log.Errorf("Failed to get QR channel: %v", err)
		}
	} else {
		go func() {
			for evt := range ch {
				if evt.Event == "code" {
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				} else {
					log.Infof("QR channel result: %s", evt.Event)
				}
			}
		}()
	}

	cli.AddEventHandler(handler)
	err = cli.Connect()
	if err != nil {
		log.Errorf("Failed to connect: %v", err)
		return
	}

	c := make(chan os.Signal)
	input := make(chan string)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer close(input)
		scan := bufio.NewScanner(os.Stdin)
		for scan.Scan() {
			line := strings.TrimSpace(scan.Text())
			if len(line) > 0 {
				input <- line
			}
		}
	}()
args := os.Args[1:]
if len(args) > 0 && args[0] != "null" {
	if len(args) > 1 && args[1] == "start" {
		tasker_type = "net.dinglisch.android.taskerm"
	} else if len(args) > 1 && args[1] == "start" {
		tasker_type = "net.dinglisch.android.tasker"
	} else {
		tasker_type = "both"
	}
	if args[0] == "both_mode" {
		is_mode = "both"
	} else if args[0] == "receive_mode" {
		is_mode = "receive" 
	} else if args[0] == "send_mode" {
		is_mode = "send"
	} else {
		handleCmd(strings.ToLower(args[0]), args[1:])
	return
	}
}
	for {
		select {
		case <-c:
			log.Infof("Interrupt received, exiting")
if is_mode == "both" || is_mode == "receive" || is_mode == "send" {
	kill_server()
}
			cli.Disconnect()
			return
		case cmd := <-input:
			if len(cmd) == 0 {
				log.Infof("Stdin closed, exiting")
if is_mode == "both" || is_mode == "receive" || is_mode == "send" {
	kill_server()
}
				cli.Disconnect()
				return
			}
			if isWaitingForPair.Load() {
				if cmd == "r" {
					pairRejectChan <- true
				} else if cmd == "a" {
					pairRejectChan <- false
				}
				continue
			}
			args := strings.Fields(cmd)
			cmd = args[0]
			args = args[1:]
			go handleCmd(strings.ToLower(cmd), args)
		}
	}
}

func parseJID(arg string) (types.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			log.Errorf("Invalid JID %s: %v", arg, err)
			return recipient, false
		} else if recipient.User == "" {
			log.Errorf("Invalid JID %s: no server specified", arg)
			return recipient, false
		}
		return recipient, true
	}
}

func handleCmd(cmd string, args []string) {
	switch cmd {
	case "reconnect":
if is_mode == "both" || is_mode == "receive" || is_mode == "send" {
	kill_server()
}
		cli.Disconnect()
		err := cli.Connect()
		if err != nil {
			log.Errorf("Failed to connect: %v", err)
		}
	case "logout":
		err := cli.Logout()
		if err != nil {
			log.Errorf("Error logging out: %v", err)
		} else {
			log.Infof("Successfully logged out")
		}
	case "appstate":
		if len(args) < 1 {
			log.Errorf("Usage: appstate <types...>")
			return
		}
		names := []appstate.WAPatchName{appstate.WAPatchName(args[0])}
		if args[0] == "all" {
			names = []appstate.WAPatchName{appstate.WAPatchRegular, appstate.WAPatchRegularHigh, appstate.WAPatchRegularLow, appstate.WAPatchCriticalUnblockLow, appstate.WAPatchCriticalBlock}
		}
		resync := len(args) > 1 && args[1] == "resync"
		for _, name := range names {
			err := cli.FetchAppState(name, resync, false)
			if err != nil {
				log.Errorf("Failed to sync app state: %v", err)
			}
		}
	case "request-appstate-key":
		if len(args) < 1 {
			log.Errorf("Usage: request-appstate-key <ids...>")
			return
		}
		var keyIDs = make([][]byte, len(args))
		for i, id := range args {
			decoded, err := hex.DecodeString(id)
			if err != nil {
				log.Errorf("Failed to decode %s as hex: %v", id, err)
				return
			}
			keyIDs[i] = decoded
		}
		cli.DangerousInternals().RequestAppStateKeys(context.Background(), keyIDs)
	case "checkuser":
		if len(args) < 1 {
			log.Errorf("Usage: checkuser <phone numbers...>")
			return
		}
		resp, err := cli.IsOnWhatsApp(args)
		if err != nil {
			log.Errorf("Failed to check if users are on WhatsApp:", err)
		} else {
			for _, item := range resp {
				if item.VerifiedName != nil {
					log.Infof("%s: on whatsapp: %t, JID: %s, business name: %s", item.Query, item.IsIn, item.JID, item.VerifiedName.Details.GetVerifiedName())
				} else {
					log.Infof("%s: on whatsapp: %t, JID: %s", item.Query, item.IsIn, item.JID)
				}
			}
		}
	case "checkupdate":
		resp, err := cli.CheckUpdate()
		if err != nil {
			log.Errorf("Failed to check for updates: %v", err)
		} else {
			log.Debugf("Version data: %#v", resp)
			if resp.ParsedVersion == store.GetWAVersion() {
				log.Infof("Client is up to date")
			} else if store.GetWAVersion().LessThan(resp.ParsedVersion) {
				log.Warnf("Client is outdated")
			} else {
				log.Infof("Client is newer than latest")
			}
		}
	case "subscribepresence":
		if len(args) < 1 {
			log.Errorf("Usage: subscribepresence <jid>")
			return
		}
		jid, ok := parseJID(args[0])
		if !ok {
			return
		}
		err := cli.SubscribePresence(jid)
		if err != nil {
			fmt.Println(err)
		}
	case "presence":
		if len(args) == 0 {
			log.Errorf("Usage: presence <available/unavailable>")
			return
		}
		fmt.Println(cli.SendPresence(types.Presence(args[0])))
	case "chatpresence":
		if len(args) == 2 {
			args = append(args, "")
		} else if len(args) < 2 {
			log.Errorf("Usage: chatpresence <jid> <composing/paused> [audio]")
			return
		}
		jid, _ := types.ParseJID(args[0])
		fmt.Println(cli.SendChatPresence(jid, types.ChatPresence(args[1]), types.ChatPresenceMedia(args[2])))
	case "privacysettings":
		resp, err := cli.TryFetchPrivacySettings(false)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%+v\n", resp)
		}
case "listusers":
		users, err := cli.Store.Contacts.GetAllContacts()
		if err != nil {
			log.Errorf("Failed to get user list: %v", err)
		} else {
			for number, user := range users {
				log.Infof("%v:%+v", number, user)
			}
		}
	case "getuser":
		if len(args) < 1 {
			log.Errorf("Usage: getuser <jids...>")
			return
		}
		var jids []types.JID
		for _, arg := range args {
			jid, ok := parseJID(arg)
			if !ok {
				return
			}
			jids = append(jids, jid)
		}
		resp, err := cli.GetUserInfo(jids)
		if err != nil {
			log.Errorf("Failed to get user info: %v", err)
		} else {
			for jid, info := range resp {
				log.Infof("%s: %+v", jid, info)
			}
		}
	case "mediaconn":
		conn, err := cli.DangerousInternals().RefreshMediaConn(false)
		if err != nil {
			log.Errorf("Failed to get media connection: %v", err)
		} else {
			log.Infof("Media connection: %+v", conn)
		}
	case "getavatar":
		if len(args) < 1 {
			log.Errorf("Usage: getavatar <jid> [existing ID] [--preview] [--community]")
			return
		}
		jid, ok := parseJID(args[0])
		if !ok {
			return
		}
		existingID := ""
		if len(args) > 2 {
			existingID = args[2]
		}
		var preview, isCommunity bool
		for _, arg := range args {
			if arg == "--preview" {
				preview = true
			} else if arg == "--community" {
				isCommunity = true
			}
		}
		pic, err := cli.GetProfilePictureInfo(jid, &whatsmeow.GetProfilePictureParams{
			Preview:     preview,
			IsCommunity: isCommunity,
			ExistingID:  existingID,
		})
		if err != nil {
			log.Errorf("Failed to get avatar: %v", err)
		} else if pic != nil {
			log.Infof("Got avatar ID %s: %s", pic.ID, pic.URL)
		} else {
			log.Infof("No avatar found")
		}
	case "getgroup":
		if len(args) < 1 {
			log.Errorf("Usage: getgroup <jid>")
			return
		}
		group, ok := parseJID(args[0])
		if !ok {
			return
		} else if group.Server != types.GroupServer {
			log.Errorf("Input must be a group JID (@%s)", types.GroupServer)
			return
		}
		resp, err := cli.GetGroupInfo(group)
		if err != nil {
			log.Errorf("Failed to get group info: %v", err)
		} else {
			log.Infof("Group info: %+v", resp)
		}
	case "subgroups":
		if len(args) < 1 {
			log.Errorf("Usage: subgroups <jid>")
			return
		}
		group, ok := parseJID(args[0])
		if !ok {
			return
		} else if group.Server != types.GroupServer {
			log.Errorf("Input must be a group JID (@%s)", types.GroupServer)
			return
		}
		resp, err := cli.GetSubGroups(group)
		if err != nil {
			log.Errorf("Failed to get subgroups: %v", err)
		} else {
			for _, sub := range resp {
				log.Infof("Subgroup: %+v", sub)
			}
		}
	case "communityparticipants":
		if len(args) < 1 {
			log.Errorf("Usage: communityparticipants <jid>")
			return
		}
		group, ok := parseJID(args[0])
		if !ok {
			return
		} else if group.Server != types.GroupServer {
			log.Errorf("Input must be a group JID (@%s)", types.GroupServer)
			return
		}
		resp, err := cli.GetLinkedGroupsParticipants(group)
		if err != nil {
			log.Errorf("Failed to get community participants: %v", err)
		} else {
			log.Infof("Community participants: %+v", resp)
		}
	case "listgroups":
		groups, err := cli.GetJoinedGroups()
		if err != nil {
			log.Errorf("Failed to get group list: %v", err)
		} else {
			for _, group := range groups {
				log.Infof("%+v", group)
			}
		}
	case "getinvitelink":
		if len(args) < 1 {
			log.Errorf("Usage: getinvitelink <jid> [--reset]")
			return
		}
		group, ok := parseJID(args[0])
		if !ok {
			return
		} else if group.Server != types.GroupServer {
			log.Errorf("Input must be a group JID (@%s)", types.GroupServer)
			return
		}
		resp, err := cli.GetGroupInviteLink(group, len(args) > 1 && args[1] == "--reset")
		if err != nil {
			log.Errorf("Failed to get group invite link: %v", err)
		} else {
			log.Infof("Group invite link: %s", resp)
		}
	case "queryinvitelink":
		if len(args) < 1 {
			log.Errorf("Usage: queryinvitelink <link>")
			return
		}
		resp, err := cli.GetGroupInfoFromLink(args[0])
		if err != nil {
			log.Errorf("Failed to resolve group invite link: %v", err)
		} else {
			log.Infof("Group info: %+v", resp)
		}
	case "querybusinesslink":
		if len(args) < 1 {
			log.Errorf("Usage: querybusinesslink <link>")
			return
		}
		resp, err := cli.ResolveBusinessMessageLink(args[0])
		if err != nil {
			log.Errorf("Failed to resolve business message link: %v", err)
		} else {
			log.Infof("Business info: %+v", resp)
		}
	case "joininvitelink":
		if len(args) < 1 {
			log.Errorf("Usage: acceptinvitelink <link>")
			return
		}
		groupID, err := cli.JoinGroupWithLink(args[0])
		if err != nil {
			log.Errorf("Failed to join group via invite link: %v", err)
		} else {
			log.Infof("Joined %s", groupID)
		}
	case "getstatusprivacy":
		resp, err := cli.GetStatusPrivacy()
		fmt.Println(err)
		fmt.Println(resp)
	case "setdisappeartimer":
		if len(args) < 2 {
			log.Errorf("Usage: setdisappeartimer <jid> <days>")
			return
		}
		days, err := strconv.Atoi(args[1])
		if err != nil {
			log.Errorf("Invalid duration: %v", err)
			return
		}
		recipient, ok := parseJID(args[0])
		if !ok {
			return
		}
		err = cli.SetDisappearingTimer(recipient, time.Duration(days)*24*time.Hour)
		if err != nil {
			log.Errorf("Failed to set disappearing timer: %v", err)
		}
	case "send":
		if len(args) < 2 {
			log.Errorf("Usage: send <jid> <text>")
			return
		}
		recipient, ok := parseJID(args[0])
		if !ok {
			return
		}
		msg := &waProto.Message{Conversation: proto.String(strings.Join(args[1:], " "))}
		resp, err := cli.SendMessage(context.Background(), recipient, msg)
		if err != nil {
			log.Errorf("Error sending message: %v", err)
		} else {
			log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
		}
	case "sendpoll":
		if len(args) < 7 {
			log.Errorf("Usage: sendpoll <jid> <max answers> <question> -- <option 1> / <option 2> / ...")
			return
		}
		recipient, ok := parseJID(args[0])
		if !ok {
			return
		}
		maxAnswers, err := strconv.Atoi(args[1])
		if err != nil {
			log.Errorf("Number of max answers must be an integer")
			return
		}
		remainingArgs := strings.Join(args[2:], " ")
		question, optionsStr, _ := strings.Cut(remainingArgs, "--")
		question = strings.TrimSpace(question)
		options := strings.Split(optionsStr, "/")
		for i, opt := range options {
			options[i] = strings.TrimSpace(opt)
		}
if is_mode == "both" || is_mode == "receive" {
			msgID := whatsmeow.GenerateMessageID()
			json_question, _ := json.Marshal(args[2])
			if tasker_type == "net.dinglisch.android.taskerm" {
				file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.taskerm/files/whatsmeow3/.tmp", msgID + "_" + args[0] + "_poll_question_")
				encoded_json_data := base64.StdEncoding.EncodeToString([]byte(json_question))
				_, _ = file.WriteString(encoded_json_data)
				_ = file.Close()
			} else if tasker_type == "net.dinglisch.android.tasker" {
				file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.tasker/files/whatsmeow3/.tmp", msgID + "_" + args[0] + "_poll_question_")
				encoded_json_data := base64.StdEncoding.EncodeToString([]byte(json_question))
				_, _ = file.WriteString(encoded_json_data)
				_ = file.Close()
			}
			resp, err := cli.SendMessage(context.Background(), recipient, cli.BuildPollCreation(question, options, maxAnswers), whatsmeow.SendRequestExtra{ID: msgID})
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
				log.Infof("Message ID : %s", msgID)
			}
			return
		}
		resp, err := cli.SendMessage(context.Background(), recipient, cli.BuildPollCreation(question, options, maxAnswers))
		if err != nil {
			log.Errorf("Error sending message: %v", err)
		} else {
			log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
		}
	case "multisend":
		if len(args) < 3 {
			log.Errorf("Usage: multisend <jids...> -- <text>")
			return
		}
		var recipients []types.JID
		for len(args) > 0 && args[0] != "--" {
			recipient, ok := parseJID(args[0])
			args = args[1:]
			if !ok {
				return
			}
			recipients = append(recipients, recipient)
		}
		if len(args) == 0 {
			log.Errorf("Usage: multisend <jids...> -- <text> (the -- is required)")
			return
		}
		msg := &waProto.Message{Conversation: proto.String(strings.Join(args[1:], " "))}
		for _, recipient := range recipients {
			go func(recipient types.JID) {
				resp, err := cli.SendMessage(context.Background(), recipient, msg)
				if err != nil {
					log.Errorf("Error sending message to %s: %v", recipient, err)
				} else {
					log.Infof("Message sent to %s (server timestamp: %s)", recipient, resp.Timestamp)
				}
			}(recipient)
		}
	case "react":
		if len(args) < 3 {
			log.Errorf("Usage: react <jid> <message ID> <reaction>")
			return
		}
		recipient, ok := parseJID(args[0])
		if !ok {
			return
		}
		messageID := args[1]
		fromMe := false
		if strings.HasPrefix(messageID, "me:") {
			fromMe = true
			messageID = messageID[len("me:"):]
		}
		reaction := args[2]
		if reaction == "remove" {
			reaction = ""
		}
		msg := &waProto.Message{
			ReactionMessage: &waProto.ReactionMessage{
				Key: &waProto.MessageKey{
					RemoteJid: proto.String(recipient.String()),
					FromMe:    proto.Bool(fromMe),
					Id:        proto.String(messageID),
				},
				Text:              proto.String(reaction),
				SenderTimestampMs: proto.Int64(time.Now().UnixMilli()),
			},
		}
		resp, err := cli.SendMessage(context.Background(), recipient, msg)
		if err != nil {
			log.Errorf("Error sending reaction: %v", err)
		} else {
			log.Infof("Reaction sent (server timestamp: %s)", resp.Timestamp)
		}
	case "revoke":
		if len(args) < 2 {
			log.Errorf("Usage: revoke <jid> <message ID>")
			return
		}
		recipient, ok := parseJID(args[0])
		if !ok {
			return
		}
		messageID := args[1]
		resp, err := cli.SendMessage(context.Background(), recipient, cli.BuildRevoke(recipient, types.EmptyJID, messageID))
		if err != nil {
			log.Errorf("Error sending revocation: %v", err)
		} else {
			log.Infof("Revocation sent (server timestamp: %s)", resp.Timestamp)
		}
	case "senddoc":
		if len(args) < 3 {
			log.Errorf("Usage: senddoc <jid> <document path> <title> [mime-type]")
			return
		}
		recipient, ok := parseJID(args[0])
		if !ok {
			return
		}
		data, err := os.ReadFile(args[1])
		if err != nil {
			log.Errorf("Failed to read %s: %v", args[0], err)
			return
		}
		uploaded, err := cli.Upload(context.Background(), data, whatsmeow.MediaDocument)
		if err != nil {
			log.Errorf("Failed to upload file: %v", err)
			return
		}
		if len(args) < 4 {
			msg := &waProto.Message{DocumentMessage: &waProto.DocumentMessage{
				Title:         proto.String(args[2]),
				Url:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String(http.DetectContentType(data)),
				FileEncSha256: uploaded.FileEncSHA256,
				FileSha256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(data))),
			}}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
		if err != nil {
			log.Errorf("Error sending document message: %v", err)
		} else {
			log.Infof("Document message sent (server timestamp: %s)", resp.Timestamp)
		}
		} else {
			msg := &waProto.Message{DocumentMessage: &waProto.DocumentMessage{
				Title:         proto.String(args[2]),
				Url:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String(args[3]),
				FileEncSha256: uploaded.FileEncSHA256,
				FileSha256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(data))),
			}}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
		    if err != nil {
			    log.Errorf("Error sending document message: %v", err)
		    } else {
			    log.Infof("Document message sent (server timestamp: %s)", resp.Timestamp)
		    }
		}
	
	case "sendvid":
		if len(args) < 2 {
			log.Errorf("Usage: sendvid <jid> <video path> [thumbnail path] [caption <note: if sending caption without specifying \"thumbnail path\", then put \"null\" in it's place.>]")
			return
		}
		recipient, ok := parseJID(args[0])
		if !ok {
			return
		}
		data, err := os.ReadFile(args[1])
		if err != nil {
			log.Errorf("Failed to read %s: %v", args[0], err)
			return
		}
		uploaded, err := cli.Upload(context.Background(), data, whatsmeow.MediaVideo)
		if err != nil {
			log.Errorf("Failed to upload file: %v", err)
			return
		}
		
		if len(args) < 3 {
			log.Errorf("No preview thumbnail specified")
			log.Infof("The video message will still be sent, but preview won't be available")
			msg := &waProto.Message{VideoMessage: &waProto.VideoMessage{
				Caption:       proto.String(strings.Join(args[2:], " ")),
				Url:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String(http.DetectContentType(data)),
				FileEncSha256: uploaded.FileEncSHA256,
				FileSha256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(data))),
			}}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending video message: %v", err)
			} else {
				log.Infof("Video message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else {
			if args[2] == "null" {
				log.Errorf("No preview thumbnail specified")
				log.Infof("The video message will still be sent, but preview won't be available")
				msg := &waProto.Message{VideoMessage: &waProto.VideoMessage{
					Caption:       proto.String(strings.Join(args[3:], " ")),
					Url:           proto.String(uploaded.URL),
					DirectPath:    proto.String(uploaded.DirectPath),
					MediaKey:      uploaded.MediaKey,
					Mimetype:      proto.String(http.DetectContentType(data)),
					FileEncSha256: uploaded.FileEncSHA256,
					FileSha256:    uploaded.FileSHA256,
					FileLength:    proto.Uint64(uint64(len(data))),
				}}
				resp, err := cli.SendMessage(context.Background(), recipient, msg)
				if err != nil {
					log.Errorf("Error sending video message: %v", err)
				} else {
					log.Infof("Video message sent (server timestamp: %s)", resp.Timestamp)
				}
			} else {
				jpegImageFile, jpegErr := os.Open(args[2])
	        	if jpegErr != nil {
		        	log.Errorf("Failed to find preview thumbnail file.")
		        	return
	        	}
				defer jpegImageFile.Close()
				jpegFileinfo, _ := jpegImageFile.Stat()
				var jpegSize int64 = jpegFileinfo.Size()
				jpegBytes := make([]byte, jpegSize)
				jpegBuffer := bufio.NewReader(jpegImageFile)
				_, jpegErr = jpegBuffer.Read(jpegBytes)
				thumbnailResp, err := cli.Upload(context.Background(), jpegBytes, whatsmeow.MediaImage)
				if err != nil {
					log.Errorf("Failed to upload preview thumbnail file: %v", err)
					return
				}
				msg := &waProto.Message{VideoMessage: &waProto.VideoMessage{
					Caption:       proto.String(strings.Join(args[3:], " ")),
					Url:           proto.String(uploaded.URL),
					DirectPath:    proto.String(uploaded.DirectPath),
					ThumbnailDirectPath: &thumbnailResp.DirectPath,
					ThumbnailSha256: thumbnailResp.FileSHA256,
					ThumbnailEncSha256: thumbnailResp.FileEncSHA256,
					JpegThumbnail: jpegBytes,
					MediaKey:      uploaded.MediaKey,
					Mimetype:      proto.String(http.DetectContentType(data)),
					FileEncSha256: uploaded.FileEncSHA256,
					FileSha256:    uploaded.FileSHA256,
					FileLength:    proto.Uint64(uint64(len(data))),
				}}
				resp, err := cli.SendMessage(context.Background(), recipient, msg)
				if err != nil {
					log.Errorf("Error sending video message: %v", err)
				} else {
					log.Infof("Video message sent (server timestamp: %s)", resp.Timestamp)
				}
			}
		}
		
	case "sendaudio":
			if len(args) < 2 {
			log.Errorf("Usage: sendaudio <jid> <audio path>")
			return
		}
		recipient, ok := parseJID(args[0])
		if !ok {
			return
		}
		data, err := os.ReadFile(args[1])
		if err != nil {
			log.Errorf("Failed to read %s: %v", args[0], err)
			return
		}
		uploaded, err := cli.Upload(context.Background(), data, whatsmeow.MediaAudio)
		if err != nil {
			log.Errorf("Failed to upload file: %v", err)
			return
		}
		msg := &waProto.Message{AudioMessage: &waProto.AudioMessage{
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String("audio/ogg; codecs=opus"),
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
		}}
		resp, err := cli.SendMessage(context.Background(), recipient, msg)
		if err != nil {
			log.Errorf("Error sending audio message: %v", err)
		} else {
			log.Infof("Audio message sent (server timestamp: %s)", resp.Timestamp)
		}
	
	case "sendimg":
			if len(args) < 2 {
			log.Errorf("Usage: sendimg <jid> <image path> [thumbnail path] [caption <note: if sending caption without specifying \"thumbnail path\", then put \"null\" in it's place.>]")
			return
		}
		recipient, ok := parseJID(args[0])
		if !ok {
			return
		}
		data, err := os.ReadFile(args[1])
		if err != nil {
			log.Errorf("Failed to read %s: %v", args[0], err)
			return
		}
		uploaded, err := cli.Upload(context.Background(), data, whatsmeow.MediaImage)
		if err != nil {
			log.Errorf("Failed to upload file: %v", err)
			return
		}
		if len(args) < 3 {
			log.Errorf("No preview thumbnail specified")
			log.Infof("The image message will still be sent, but preview won't be available")
			msg := &waProto.Message{ImageMessage: &waProto.ImageMessage{
				Caption:       proto.String(strings.Join(args[2:], " ")),
				Url:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String(http.DetectContentType(data)),
				FileEncSha256: uploaded.FileEncSHA256,
				FileSha256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(data))),
			}}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending image message: %v", err)
			} else {
				log.Infof("Image message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else {
		    if args[2] == "null" {
		        log.Errorf("No preview thumbnail specified")
		        log.Infof("The image message will still be sent, but preview won't be available")
			    msg := &waProto.Message{ImageMessage: &waProto.ImageMessage{
				    Caption:       proto.String(strings.Join(args[3:], " ")),
				    Url:           proto.String(uploaded.URL),
				    DirectPath:    proto.String(uploaded.DirectPath),
				    MediaKey:      uploaded.MediaKey,
				    Mimetype:      proto.String(http.DetectContentType(data)),
				    FileEncSha256: uploaded.FileEncSHA256,
				    FileSha256:    uploaded.FileSHA256,
				    FileLength:    proto.Uint64(uint64(len(data))),
			    }}
			    resp, err := cli.SendMessage(context.Background(), recipient, msg)
			    if err != nil {
				    log.Errorf("Error sending image message: %v", err)
			    } else {
				    log.Infof("Image message sent (server timestamp: %s)", resp.Timestamp)
			    }     
		    } else {
				jpegImageFile, jpegErr := os.Open(args[2])
				if jpegErr != nil {
				log.Errorf("Failed to find preview thumbnail file.")
				return
				}
				defer jpegImageFile.Close()
				jpegFileinfo, _ := jpegImageFile.Stat()
				var jpegSize int64 = jpegFileinfo.Size()
				jpegBytes := make([]byte, jpegSize)
				jpegBuffer := bufio.NewReader(jpegImageFile)
				_, jpegErr = jpegBuffer.Read(jpegBytes)
				thumbnailResp, err := cli.Upload(context.Background(), jpegBytes, whatsmeow.MediaImage)
				if err != nil {
					log.Errorf("Failed to upload preview thumbnail file: %v", err)
					return
				}
			    msg := &waProto.Message{ImageMessage: &waProto.ImageMessage{
				    Caption:       proto.String(strings.Join(args[3:], " ")),
				    Url:           proto.String(uploaded.URL),
				    DirectPath:    proto.String(uploaded.DirectPath),
				    ThumbnailDirectPath: &thumbnailResp.DirectPath,
				    ThumbnailSha256: thumbnailResp.FileSHA256,
				    ThumbnailEncSha256: thumbnailResp.FileEncSHA256,
				    JpegThumbnail: jpegBytes,
				    MediaKey:      uploaded.MediaKey,
				    Mimetype:      proto.String(http.DetectContentType(data)),
				    FileEncSha256: uploaded.FileEncSHA256,
				    FileSha256:    uploaded.FileSHA256,
				    FileLength:    proto.Uint64(uint64(len(data))),
			    }}
			    resp, err := cli.SendMessage(context.Background(), recipient, msg)
			    if err != nil {
				    log.Errorf("Error sending image message: %v", err)
			    } else {
				    log.Infof("Image message sent (server timestamp: %s)", resp.Timestamp)
			    }
		    }
			
		}
		
	case "sendbutton":
		if len(args) < 5 {
			log.Errorf("Usage: sendbutton <jid> <title> <text body> <footer> <button1> [button2] [button3] (Note: [] is optional)")
			return
		}
		recipient, ok := parseJID(args[0])
		if !ok {
			return
		}
		if len(args) == 5 || (len(args) > 5 && args[5] == "") {
			msg := &waProto.Message{
				ViewOnceMessage: &waProto.FutureProofMessage{
					Message: &waProto.Message{
						ButtonsMessage: &waProto.ButtonsMessage{
							HeaderType: waProto.ButtonsMessage_TEXT.Enum(),
							Header: &waProto.ButtonsMessage_Text{
								Text: args[1], 
							},
							ContentText: proto.String(args[2]),
							FooterText:  proto.String(args[3]),
							Buttons: []*waProto.ButtonsMessage_Button{
								{
									ButtonId: proto.String("button1"),
									ButtonText: &waProto.ButtonsMessage_Button_ButtonText{
										DisplayText: proto.String(args[4]),
									},
									Type:           waProto.ButtonsMessage_Button_RESPONSE.Enum(),
									NativeFlowInfo: &waProto.ButtonsMessage_Button_NativeFlowInfo{},
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else if len(args) == 6 || (len(args) > 6 && args[6] == "") {
			msg := &waProto.Message{
				ViewOnceMessage: &waProto.FutureProofMessage{
					Message: &waProto.Message{
						ButtonsMessage: &waProto.ButtonsMessage{
							HeaderType: waProto.ButtonsMessage_TEXT.Enum(),
							Header: &waProto.ButtonsMessage_Text{
								Text: args[1], 
							},
							ContentText: proto.String(args[2]),
							FooterText:  proto.String(args[3]),
							Buttons: []*waProto.ButtonsMessage_Button{
								{
									ButtonId: proto.String("button1"),
									ButtonText: &waProto.ButtonsMessage_Button_ButtonText{
										DisplayText: proto.String(args[4]),
									},
									Type:           waProto.ButtonsMessage_Button_RESPONSE.Enum(),
									NativeFlowInfo: &waProto.ButtonsMessage_Button_NativeFlowInfo{},
								},
								{
									ButtonId: proto.String("button2"),
									ButtonText: &waProto.ButtonsMessage_Button_ButtonText{
										DisplayText: proto.String(args[5]),
									},
									Type:           waProto.ButtonsMessage_Button_RESPONSE.Enum(),
									NativeFlowInfo: &waProto.ButtonsMessage_Button_NativeFlowInfo{},
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else if len(args) == 7 || (len(args) > 7 && args[7] == "") {
			msg := &waProto.Message{
				ViewOnceMessage: &waProto.FutureProofMessage{
					Message: &waProto.Message{
						ButtonsMessage: &waProto.ButtonsMessage{
							HeaderType: waProto.ButtonsMessage_TEXT.Enum(),
							Header: &waProto.ButtonsMessage_Text{
								Text: args[1], 
							},
							ContentText: proto.String(args[2]),
							FooterText:  proto.String(args[3]),
							Buttons: []*waProto.ButtonsMessage_Button{
								{
									ButtonId: proto.String("button1"),
									ButtonText: &waProto.ButtonsMessage_Button_ButtonText{
										DisplayText: proto.String(args[4]),
									},
									Type:           waProto.ButtonsMessage_Button_RESPONSE.Enum(),
									NativeFlowInfo: &waProto.ButtonsMessage_Button_NativeFlowInfo{},
								},
								{
									ButtonId: proto.String("button2"),
									ButtonText: &waProto.ButtonsMessage_Button_ButtonText{
										DisplayText: proto.String(args[5]),
									},
									Type:           waProto.ButtonsMessage_Button_RESPONSE.Enum(),
									NativeFlowInfo: &waProto.ButtonsMessage_Button_NativeFlowInfo{},
								},
								{
									ButtonId: proto.String("button3"),
									ButtonText: &waProto.ButtonsMessage_Button_ButtonText{
										DisplayText: proto.String(args[6]),
									},
									Type:           waProto.ButtonsMessage_Button_RESPONSE.Enum(),
									NativeFlowInfo: &waProto.ButtonsMessage_Button_NativeFlowInfo{},
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		}
		
	case "sendlist":
		if len(args) < 8 {
			log.Errorf("Usage: sendlist <jid> <title> <text body> <footer> <button text> <list header> <list title 1> <list description 1> [list title X] [list description X] (Note: Upto 15 item pairs. [] is optional)")
			return
		}
		recipient, ok := parseJID(args[0])
		if !ok {
			return
		}
		if len(args) == 8 || (len(args) > 8 && args[8] == "") {
			msg := &waProto.Message{
				ListMessage: &waProto.ListMessage{
					Title:       proto.String(args[1]),
					Description: proto.String(args[2]),
					FooterText:  proto.String(args[3]),
					ButtonText:  proto.String(args[4]),
					ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
					Sections: []*waProto.ListMessage_Section{
						{
							Title: proto.String(args[5]),
							Rows: []*waProto.ListMessage_Row{
								{
									RowId:       proto.String("id1"),
									Title:       proto.String(args[6]),
									Description: proto.String(args[7]),
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else if len(args) == 10 || (len(args) > 10 && args[10] == "") {
			msg := &waProto.Message{
				ListMessage: &waProto.ListMessage{
					Title:       proto.String(args[1]),
					Description: proto.String(args[2]),
					FooterText:  proto.String(args[3]),
					ButtonText:  proto.String(args[4]),
					ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
					Sections: []*waProto.ListMessage_Section{
						{
							Title: proto.String(args[5]),
							Rows: []*waProto.ListMessage_Row{
								{
									RowId:       proto.String("id1"),
									Title:       proto.String(args[6]),
									Description: proto.String(args[7]),
								},
								{
									RowId:       proto.String("id2"),
									Title:       proto.String(args[8]),
									Description: proto.String(args[9]),
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else if len(args) == 12 || (len(args) > 12 && args[12] == "") {
			msg := &waProto.Message{
				ListMessage: &waProto.ListMessage{
					Title:       proto.String(args[1]),
					Description: proto.String(args[2]),
					FooterText:  proto.String(args[3]),
					ButtonText:  proto.String(args[4]),
					ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
					Sections: []*waProto.ListMessage_Section{
						{
							Title: proto.String(args[5]),
							Rows: []*waProto.ListMessage_Row{
								{
									RowId:       proto.String("id1"),
									Title:       proto.String(args[6]),
									Description: proto.String(args[7]),
								},
								{
									RowId:       proto.String("id2"),
									Title:       proto.String(args[8]),
									Description: proto.String(args[9]),
								},
								{
									RowId:       proto.String("id3"),
									Title:       proto.String(args[10]),
									Description: proto.String(args[11]),
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else if len(args) == 14 || (len(args) > 14 && args[14] == "") {
			msg := &waProto.Message{
				ListMessage: &waProto.ListMessage{
					Title:       proto.String(args[1]),
					Description: proto.String(args[2]),
					FooterText:  proto.String(args[3]),
					ButtonText:  proto.String(args[4]),
					ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
					Sections: []*waProto.ListMessage_Section{
						{
							Title: proto.String(args[5]),
							Rows: []*waProto.ListMessage_Row{
								{
									RowId:       proto.String("id1"),
									Title:       proto.String(args[6]),
									Description: proto.String(args[7]),
								},
								{
									RowId:       proto.String("id2"),
									Title:       proto.String(args[8]),
									Description: proto.String(args[9]),
								},
								{
									RowId:       proto.String("id3"),
									Title:       proto.String(args[10]),
									Description: proto.String(args[11]),
								},
								{
									RowId:       proto.String("id4"),
									Title:       proto.String(args[12]),
									Description: proto.String(args[13]),
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else if len(args) == 16 || (len(args) > 16 && args[16] == "") {
			msg := &waProto.Message{
				ListMessage: &waProto.ListMessage{
					Title:       proto.String(args[1]),
					Description: proto.String(args[2]),
					FooterText:  proto.String(args[3]),
					ButtonText:  proto.String(args[4]),
					ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
					Sections: []*waProto.ListMessage_Section{
						{
							Title: proto.String(args[5]),
							Rows: []*waProto.ListMessage_Row{
								{
									RowId:       proto.String("id1"),
									Title:       proto.String(args[6]),
									Description: proto.String(args[7]),
								},
								{
									RowId:       proto.String("id2"),
									Title:       proto.String(args[8]),
									Description: proto.String(args[9]),
								},
								{
									RowId:       proto.String("id3"),
									Title:       proto.String(args[10]),
									Description: proto.String(args[11]),
								},
								{
									RowId:       proto.String("id4"),
									Title:       proto.String(args[12]),
									Description: proto.String(args[13]),
								},
								{
									RowId:       proto.String("id5"),
									Title:       proto.String(args[14]),
									Description: proto.String(args[15]),
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else if len(args) == 18 || (len(args) > 18 && args[18] == "") {
			msg := &waProto.Message{
				ListMessage: &waProto.ListMessage{
					Title:       proto.String(args[1]),
					Description: proto.String(args[2]),
					FooterText:  proto.String(args[3]),
					ButtonText:  proto.String(args[4]),
					ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
					Sections: []*waProto.ListMessage_Section{
						{
							Title: proto.String(args[5]),
							Rows: []*waProto.ListMessage_Row{
								{
									RowId:       proto.String("id1"),
									Title:       proto.String(args[6]),
									Description: proto.String(args[7]),
								},
								{
									RowId:       proto.String("id2"),
									Title:       proto.String(args[8]),
									Description: proto.String(args[9]),
								},
								{
									RowId:       proto.String("id3"),
									Title:       proto.String(args[10]),
									Description: proto.String(args[11]),
								},
								{
									RowId:       proto.String("id4"),
									Title:       proto.String(args[12]),
									Description: proto.String(args[13]),
								},
								{
									RowId:       proto.String("id5"),
									Title:       proto.String(args[14]),
									Description: proto.String(args[15]),
								},
								{
									RowId:       proto.String("id6"),
									Title:       proto.String(args[16]),
									Description: proto.String(args[17]),
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else if len(args) == 20 || (len(args) > 20 && args[20] == "") {
			msg := &waProto.Message{
				ListMessage: &waProto.ListMessage{
					Title:       proto.String(args[1]),
					Description: proto.String(args[2]),
					FooterText:  proto.String(args[3]),
					ButtonText:  proto.String(args[4]),
					ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
					Sections: []*waProto.ListMessage_Section{
						{
							Title: proto.String(args[5]),
							Rows: []*waProto.ListMessage_Row{
								{
									RowId:       proto.String("id1"),
									Title:       proto.String(args[6]),
									Description: proto.String(args[7]),
								},
								{
									RowId:       proto.String("id2"),
									Title:       proto.String(args[8]),
									Description: proto.String(args[9]),
								},
								{
									RowId:       proto.String("id3"),
									Title:       proto.String(args[10]),
									Description: proto.String(args[11]),
								},
								{
									RowId:       proto.String("id4"),
									Title:       proto.String(args[12]),
									Description: proto.String(args[13]),
								},
								{
									RowId:       proto.String("id5"),
									Title:       proto.String(args[14]),
									Description: proto.String(args[15]),
								},
								{
									RowId:       proto.String("id6"),
									Title:       proto.String(args[16]),
									Description: proto.String(args[17]),
								},
								{
									RowId:       proto.String("id7"),
									Title:       proto.String(args[18]),
									Description: proto.String(args[19]),
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else if len(args) == 22 || (len(args) > 22 && args[22] == "") {
			msg := &waProto.Message{
				ListMessage: &waProto.ListMessage{
					Title:       proto.String(args[1]),
					Description: proto.String(args[2]),
					FooterText:  proto.String(args[3]),
					ButtonText:  proto.String(args[4]),
					ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
					Sections: []*waProto.ListMessage_Section{
						{
							Title: proto.String(args[5]),
							Rows: []*waProto.ListMessage_Row{
								{
									RowId:       proto.String("id1"),
									Title:       proto.String(args[6]),
									Description: proto.String(args[7]),
								},
								{
									RowId:       proto.String("id2"),
									Title:       proto.String(args[8]),
									Description: proto.String(args[9]),
								},
								{
									RowId:       proto.String("id3"),
									Title:       proto.String(args[10]),
									Description: proto.String(args[11]),
								},
								{
									RowId:       proto.String("id4"),
									Title:       proto.String(args[12]),
									Description: proto.String(args[13]),
								},
								{
									RowId:       proto.String("id5"),
									Title:       proto.String(args[14]),
									Description: proto.String(args[15]),
								},
								{
									RowId:       proto.String("id6"),
									Title:       proto.String(args[16]),
									Description: proto.String(args[17]),
								},
								{
									RowId:       proto.String("id7"),
									Title:       proto.String(args[18]),
									Description: proto.String(args[19]),
								},
								{
									RowId:       proto.String("id8"),
									Title:       proto.String(args[20]),
									Description: proto.String(args[21]),
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else if len(args) == 24 || (len(args) > 24 && args[24] == "") {
			msg := &waProto.Message{
				ListMessage: &waProto.ListMessage{
					Title:       proto.String(args[1]),
					Description: proto.String(args[2]),
					FooterText:  proto.String(args[3]),
					ButtonText:  proto.String(args[4]),
					ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
					Sections: []*waProto.ListMessage_Section{
						{
							Title: proto.String(args[5]),
							Rows: []*waProto.ListMessage_Row{
								{
									RowId:       proto.String("id1"),
									Title:       proto.String(args[6]),
									Description: proto.String(args[7]),
								},
								{
									RowId:       proto.String("id2"),
									Title:       proto.String(args[8]),
									Description: proto.String(args[9]),
								},
								{
									RowId:       proto.String("id3"),
									Title:       proto.String(args[10]),
									Description: proto.String(args[11]),
								},
								{
									RowId:       proto.String("id4"),
									Title:       proto.String(args[12]),
									Description: proto.String(args[13]),
								},
								{
									RowId:       proto.String("id5"),
									Title:       proto.String(args[14]),
									Description: proto.String(args[15]),
								},
								{
									RowId:       proto.String("id6"),
									Title:       proto.String(args[16]),
									Description: proto.String(args[17]),
								},
								{
									RowId:       proto.String("id7"),
									Title:       proto.String(args[18]),
									Description: proto.String(args[19]),
								},
								{
									RowId:       proto.String("id8"),
									Title:       proto.String(args[20]),
									Description: proto.String(args[21]),
								},
								{
									RowId:       proto.String("id9"),
									Title:       proto.String(args[22]),
									Description: proto.String(args[23]),
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else if len(args) == 26 || (len(args) > 26 && args[26] == "") {
			msg := &waProto.Message{
				ListMessage: &waProto.ListMessage{
					Title:       proto.String(args[1]),
					Description: proto.String(args[2]),
					FooterText:  proto.String(args[3]),
					ButtonText:  proto.String(args[4]),
					ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
					Sections: []*waProto.ListMessage_Section{
						{
							Title: proto.String(args[5]),
							Rows: []*waProto.ListMessage_Row{
								{
									RowId:       proto.String("id1"),
									Title:       proto.String(args[6]),
									Description: proto.String(args[7]),
								},
								{
									RowId:       proto.String("id2"),
									Title:       proto.String(args[8]),
									Description: proto.String(args[9]),
								},
								{
									RowId:       proto.String("id3"),
									Title:       proto.String(args[10]),
									Description: proto.String(args[11]),
								},
								{
									RowId:       proto.String("id4"),
									Title:       proto.String(args[12]),
									Description: proto.String(args[13]),
								},
								{
									RowId:       proto.String("id5"),
									Title:       proto.String(args[14]),
									Description: proto.String(args[15]),
								},
								{
									RowId:       proto.String("id6"),
									Title:       proto.String(args[16]),
									Description: proto.String(args[17]),
								},
								{
									RowId:       proto.String("id7"),
									Title:       proto.String(args[18]),
									Description: proto.String(args[19]),
								},
								{
									RowId:       proto.String("id8"),
									Title:       proto.String(args[20]),
									Description: proto.String(args[21]),
								},
								{
									RowId:       proto.String("id9"),
									Title:       proto.String(args[22]),
									Description: proto.String(args[23]),
								},
								{
									RowId:       proto.String("id10"),
									Title:       proto.String(args[24]),
									Description: proto.String(args[25]),
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else if len(args) == 28 || (len(args) > 28 && args[28] == "") {
			msg := &waProto.Message{
				ListMessage: &waProto.ListMessage{
					Title:       proto.String(args[1]),
					Description: proto.String(args[2]),
					FooterText:  proto.String(args[3]),
					ButtonText:  proto.String(args[4]),
					ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
					Sections: []*waProto.ListMessage_Section{
						{
							Title: proto.String(args[5]),
							Rows: []*waProto.ListMessage_Row{
								{
									RowId:       proto.String("id1"),
									Title:       proto.String(args[6]),
									Description: proto.String(args[7]),
								},
								{
									RowId:       proto.String("id2"),
									Title:       proto.String(args[8]),
									Description: proto.String(args[9]),
								},
								{
									RowId:       proto.String("id3"),
									Title:       proto.String(args[10]),
									Description: proto.String(args[11]),
								},
								{
									RowId:       proto.String("id4"),
									Title:       proto.String(args[12]),
									Description: proto.String(args[13]),
								},
								{
									RowId:       proto.String("id5"),
									Title:       proto.String(args[14]),
									Description: proto.String(args[15]),
								},
								{
									RowId:       proto.String("id6"),
									Title:       proto.String(args[16]),
									Description: proto.String(args[17]),
								},
								{
									RowId:       proto.String("id7"),
									Title:       proto.String(args[18]),
									Description: proto.String(args[19]),
								},
								{
									RowId:       proto.String("id8"),
									Title:       proto.String(args[20]),
									Description: proto.String(args[21]),
								},
								{
									RowId:       proto.String("id9"),
									Title:       proto.String(args[22]),
									Description: proto.String(args[23]),
								},
								{
									RowId:       proto.String("id10"),
									Title:       proto.String(args[24]),
									Description: proto.String(args[25]),
								},
								{
									RowId:       proto.String("id11"),
									Title:       proto.String(args[26]),
									Description: proto.String(args[27]),
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else if len(args) == 30 || (len(args) > 30 && args[30] == "") {
			msg := &waProto.Message{
				ListMessage: &waProto.ListMessage{
					Title:       proto.String(args[1]),
					Description: proto.String(args[2]),
					FooterText:  proto.String(args[3]),
					ButtonText:  proto.String(args[4]),
					ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
					Sections: []*waProto.ListMessage_Section{
						{
							Title: proto.String(args[5]),
							Rows: []*waProto.ListMessage_Row{
								{
									RowId:       proto.String("id1"),
									Title:       proto.String(args[6]),
									Description: proto.String(args[7]),
								},
								{
									RowId:       proto.String("id2"),
									Title:       proto.String(args[8]),
									Description: proto.String(args[9]),
								},
								{
									RowId:       proto.String("id3"),
									Title:       proto.String(args[10]),
									Description: proto.String(args[11]),
								},
								{
									RowId:       proto.String("id4"),
									Title:       proto.String(args[12]),
									Description: proto.String(args[13]),
								},
								{
									RowId:       proto.String("id5"),
									Title:       proto.String(args[14]),
									Description: proto.String(args[15]),
								},
								{
									RowId:       proto.String("id6"),
									Title:       proto.String(args[16]),
									Description: proto.String(args[17]),
								},
								{
									RowId:       proto.String("id7"),
									Title:       proto.String(args[18]),
									Description: proto.String(args[19]),
								},
								{
									RowId:       proto.String("id8"),
									Title:       proto.String(args[20]),
									Description: proto.String(args[21]),
								},
								{
									RowId:       proto.String("id9"),
									Title:       proto.String(args[22]),
									Description: proto.String(args[23]),
								},
								{
									RowId:       proto.String("id10"),
									Title:       proto.String(args[24]),
									Description: proto.String(args[25]),
								},
								{
									RowId:       proto.String("id11"),
									Title:       proto.String(args[26]),
									Description: proto.String(args[27]),
								},
								{
									RowId:       proto.String("id12"),
									Title:       proto.String(args[28]),
									Description: proto.String(args[29]),
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else if len(args) == 32 || (len(args) > 32 && args[32] == "") {
			msg := &waProto.Message{
				ListMessage: &waProto.ListMessage{
					Title:       proto.String(args[1]),
					Description: proto.String(args[2]),
					FooterText:  proto.String(args[3]),
					ButtonText:  proto.String(args[4]),
					ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
					Sections: []*waProto.ListMessage_Section{
						{
							Title: proto.String(args[5]),
							Rows: []*waProto.ListMessage_Row{
								{
									RowId:       proto.String("id1"),
									Title:       proto.String(args[6]),
									Description: proto.String(args[7]),
								},
								{
									RowId:       proto.String("id2"),
									Title:       proto.String(args[8]),
									Description: proto.String(args[9]),
								},
								{
									RowId:       proto.String("id3"),
									Title:       proto.String(args[10]),
									Description: proto.String(args[11]),
								},
								{
									RowId:       proto.String("id4"),
									Title:       proto.String(args[12]),
									Description: proto.String(args[13]),
								},
								{
									RowId:       proto.String("id5"),
									Title:       proto.String(args[14]),
									Description: proto.String(args[15]),
								},
								{
									RowId:       proto.String("id6"),
									Title:       proto.String(args[16]),
									Description: proto.String(args[17]),
								},
								{
									RowId:       proto.String("id7"),
									Title:       proto.String(args[18]),
									Description: proto.String(args[19]),
								},
								{
									RowId:       proto.String("id8"),
									Title:       proto.String(args[20]),
									Description: proto.String(args[21]),
								},
								{
									RowId:       proto.String("id9"),
									Title:       proto.String(args[22]),
									Description: proto.String(args[23]),
								},
								{
									RowId:       proto.String("id10"),
									Title:       proto.String(args[24]),
									Description: proto.String(args[25]),
								},
								{
									RowId:       proto.String("id11"),
									Title:       proto.String(args[26]),
									Description: proto.String(args[27]),
								},
								{
									RowId:       proto.String("id12"),
									Title:       proto.String(args[28]),
									Description: proto.String(args[29]),
								},
								{
									RowId:       proto.String("id13"),
									Title:       proto.String(args[30]),
									Description: proto.String(args[31]),
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else if len(args) == 34 || (len(args) > 34 && args[34] == "") {
			msg := &waProto.Message{
				ListMessage: &waProto.ListMessage{
					Title:       proto.String(args[1]),
					Description: proto.String(args[2]),
					FooterText:  proto.String(args[3]),
					ButtonText:  proto.String(args[4]),
					ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
					Sections: []*waProto.ListMessage_Section{
						{
							Title: proto.String(args[5]),
							Rows: []*waProto.ListMessage_Row{
								{
									RowId:       proto.String("id1"),
									Title:       proto.String(args[6]),
									Description: proto.String(args[7]),
								},
								{
									RowId:       proto.String("id2"),
									Title:       proto.String(args[8]),
									Description: proto.String(args[9]),
								},
								{
									RowId:       proto.String("id3"),
									Title:       proto.String(args[10]),
									Description: proto.String(args[11]),
								},
								{
									RowId:       proto.String("id4"),
									Title:       proto.String(args[12]),
									Description: proto.String(args[13]),
								},
								{
									RowId:       proto.String("id5"),
									Title:       proto.String(args[14]),
									Description: proto.String(args[15]),
								},
								{
									RowId:       proto.String("id6"),
									Title:       proto.String(args[16]),
									Description: proto.String(args[17]),
								},
								{
									RowId:       proto.String("id7"),
									Title:       proto.String(args[18]),
									Description: proto.String(args[19]),
								},
								{
									RowId:       proto.String("id8"),
									Title:       proto.String(args[20]),
									Description: proto.String(args[21]),
								},
								{
									RowId:       proto.String("id9"),
									Title:       proto.String(args[22]),
									Description: proto.String(args[23]),
								},
								{
									RowId:       proto.String("id10"),
									Title:       proto.String(args[24]),
									Description: proto.String(args[25]),
								},
								{
									RowId:       proto.String("id11"),
									Title:       proto.String(args[26]),
									Description: proto.String(args[27]),
								},
								{
									RowId:       proto.String("id12"),
									Title:       proto.String(args[28]),
									Description: proto.String(args[29]),
								},
								{
									RowId:       proto.String("id13"),
									Title:       proto.String(args[30]),
									Description: proto.String(args[31]),
								},
								{
									RowId:       proto.String("id14"),
									Title:       proto.String(args[32]),
									Description: proto.String(args[33]),
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else if len(args) == 36 || (len(args) > 36 && args[36] == "") {
			msg := &waProto.Message{
				ListMessage: &waProto.ListMessage{
					Title:       proto.String(args[1]),
					Description: proto.String(args[2]),
					FooterText:  proto.String(args[3]),
					ButtonText:  proto.String(args[4]),
					ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
					Sections: []*waProto.ListMessage_Section{
						{
							Title: proto.String(args[5]),
							Rows: []*waProto.ListMessage_Row{
								{
									RowId:       proto.String("id1"),
									Title:       proto.String(args[6]),
									Description: proto.String(args[7]),
								},
								{
									RowId:       proto.String("id2"),
									Title:       proto.String(args[8]),
									Description: proto.String(args[9]),
								},
								{
									RowId:       proto.String("id3"),
									Title:       proto.String(args[10]),
									Description: proto.String(args[11]),
								},
								{
									RowId:       proto.String("id4"),
									Title:       proto.String(args[12]),
									Description: proto.String(args[13]),
								},
								{
									RowId:       proto.String("id5"),
									Title:       proto.String(args[14]),
									Description: proto.String(args[15]),
								},
								{
									RowId:       proto.String("id6"),
									Title:       proto.String(args[16]),
									Description: proto.String(args[17]),
								},
								{
									RowId:       proto.String("id7"),
									Title:       proto.String(args[18]),
									Description: proto.String(args[19]),
								},
								{
									RowId:       proto.String("id8"),
									Title:       proto.String(args[20]),
									Description: proto.String(args[21]),
								},
								{
									RowId:       proto.String("id9"),
									Title:       proto.String(args[22]),
									Description: proto.String(args[23]),
								},
								{
									RowId:       proto.String("id10"),
									Title:       proto.String(args[24]),
									Description: proto.String(args[25]),
								},
								{
									RowId:       proto.String("id11"),
									Title:       proto.String(args[26]),
									Description: proto.String(args[27]),
								},
								{
									RowId:       proto.String("id12"),
									Title:       proto.String(args[28]),
									Description: proto.String(args[29]),
								},
								{
									RowId:       proto.String("id13"),
									Title:       proto.String(args[30]),
									Description: proto.String(args[31]),
								},
								{
									RowId:       proto.String("id14"),
									Title:       proto.String(args[32]),
									Description: proto.String(args[33]),
								},
								{
									RowId:       proto.String("id15"),
									Title:       proto.String(args[34]),
									Description: proto.String(args[35]),
								},
							},
						},
					},
				},
			}
			resp, err := cli.SendMessage(context.Background(), recipient, msg)
			if err != nil {
				log.Errorf("Error sending message: %v", err)
			} else {
				log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
			}
		} else {
			log.Errorf("The number of <list title> does not match <list description>. Or beyond 15 item pairs")
		}
		
	case "markread":
		if len(args) < 2 {
			log.Errorf("Usage: markread <jid> <message ID 1> [message ID X] (Note: Can add multiple message IDs to mark as read. [] is optional)")
			return
		}
		recipient, ok := parseJID(args[0])
		if !ok {
			return
		}
		messageID := []string{}
		for _, message_id := range args[1:] {
			if message_id != "" {
				messageID = append(messageID, message_id)
			}
		}
		err := cli.MarkRead(messageID, time.Now(), recipient, types.EmptyJID)
		if err != nil {
			log.Errorf("Error sending mark as read: %v", err)
		} else {
			log.Infof("Mark as read sent")
		}
		
	case "setstatus":
		if len(args) == 0 {
			log.Errorf("Usage: setstatus <message>")
			return
		}
		err := cli.SetStatusMessage(strings.Join(args, " "))
		if err != nil {
			log.Errorf("Error setting status message: %v", err)
		} else {
			log.Infof("Status updated")
		}
	case "archive":
		if len(args) < 2 {
			log.Errorf("Usage: archive <jid> <action>")
			return
		}
		target, ok := parseJID(args[0])
		if !ok {
			return
		}
		action, err := strconv.ParseBool(args[1])
		if err != nil {
			log.Errorf("invalid second argument: %v", err)
			return
		}

		err = cli.SendAppState(appstate.BuildArchive(target, action, time.Time{}, nil))
		if err != nil {
			log.Errorf("Error changing chat's archive state: %v", err)
		}
	case "mute":
		if len(args) < 2 {
			log.Errorf("Usage: mute <jid> <action>")
			return
		}
		target, ok := parseJID(args[0])
		if !ok {
			return
		}
		action, err := strconv.ParseBool(args[1])
		if err != nil {
			log.Errorf("invalid second argument: %v", err)
			return
		}

		err = cli.SendAppState(appstate.BuildMute(target, action, 1*time.Hour))
		if err != nil {
			log.Errorf("Error changing chat's mute state: %v", err)
		}
	case "pin":
		if len(args) < 2 {
			log.Errorf("Usage: pin <jid> <action>")
			return
		}
		target, ok := parseJID(args[0])
		if !ok {
			return
		}
		action, err := strconv.ParseBool(args[1])
		if err != nil {
			log.Errorf("invalid second argument: %v", err)
			return
		}

		err = cli.SendAppState(appstate.BuildPin(target, action))
		if err != nil {
			log.Errorf("Error changing chat's pin state: %v", err)
		}
	}
}

var historySyncID int32
var startupTime = time.Now().Unix()

func handler(rawEvt interface{}) {
	switch evt := rawEvt.(type) {
	case *events.AppStateSyncComplete:
		if len(cli.Store.PushName) > 0 && evt.Name == appstate.WAPatchCriticalBlock {
			err := cli.SendPresence(types.PresenceAvailable)
			if err != nil {
				log.Warnf("Failed to send available presence: %v", err)
			} else {
				log.Infof("Marked self as available")
if is_mode == "both" {
					log.Infof("Receive/Send Mode Enabled")
		            log.Infof("Will Now Receive/Send Messages")
	            	go MdtestStart()
            	} else if is_mode == "receive" {
            		log.Infof("Receive Mode Enabled")
		            log.Infof("Will Now Receive Messages")
	            	go MdtestStart()
	            } else if is_mode == "send" {
	            	log.Infof("Send Mode Enabled")
		            log.Infof("Can Now Send Messages")
	            	go MdtestStart()
            	}
			}
		}
	case *events.Connected, *events.PushNameSetting:
		if len(cli.Store.PushName) == 0 {
			return
		}
		// Send presence available when connecting and when the pushname is changed.
		// This makes sure that outgoing messages always have the right pushname.
		err := cli.SendPresence(types.PresenceAvailable)
		if err != nil {
			log.Warnf("Failed to send available presence: %v", err)
		} else {
			log.Infof("Marked self as available")
if is_mode == "both" {
					log.Infof("Receive/Send Mode Enabled")
		            log.Infof("Will Now Receive/Send Messages")
	            	go MdtestStart()
            	} else if is_mode == "receive" {
            		log.Infof("Receive Mode Enabled")
		            log.Infof("Will Now Receive Messages")
	            	go MdtestStart()
	            } else if is_mode == "send" {
	            	log.Infof("Send Mode Enabled")
		            log.Infof("Can Now Send Messages")
	            	go MdtestStart()
            	}
		}
	case *events.StreamReplaced:
		os.Exit(0)
	case *events.Message:
		metaParts := []string{fmt.Sprintf("pushname: %s", evt.Info.PushName), fmt.Sprintf("timestamp: %s", evt.Info.Timestamp)}
		if evt.Info.Type != "" {
			metaParts = append(metaParts, fmt.Sprintf("type: %s", evt.Info.Type))
		}
		if evt.Info.Category != "" {
			metaParts = append(metaParts, fmt.Sprintf("category: %s", evt.Info.Category))
		}
		if evt.IsViewOnce {
			metaParts = append(metaParts, "view once")
		}
		if evt.IsViewOnce {
			metaParts = append(metaParts, "ephemeral")
		}
		if evt.IsViewOnceV2 {
			metaParts = append(metaParts, "ephemeral (v2)")
		}
		if evt.IsDocumentWithCaption {
			metaParts = append(metaParts, "document with caption")
		}
		if evt.IsEdit {
			metaParts = append(metaParts, "edit")
		}

		log.Infof("Received message %s from %s (%s): %+v", evt.Info.ID, evt.Info.SourceString(), strings.Join(metaParts, ", "), evt.Message)
go func() {
			if is_mode == "both" || is_mode == "receive" {
				message_id := fmt.Sprintf("%s", evt.Info.ID)
				sender_pushname := fmt.Sprintf("%s", evt.Info.PushName)
				sender_jid := fmt.Sprintf("%s", evt.Info.Sender)
				receiver_jid := fmt.Sprintf("%s", evt.Info.Chat)
				time_stamp := fmt.Sprintf("%s", evt.Info.Timestamp)
				json_sender_pushname, _ := json.Marshal(sender_pushname)
				json_time_stamp, _ := json.Marshal(time_stamp)
				var is_from_myself string
				if evt.Info.MessageSource.IsFromMe {
					is_from_myself = "1"
				} else {
					is_from_myself = "0"
				}
				var is_group string
				if evt.Info.MessageSource.IsGroup {
					is_group = "1"
				} else {
					is_group = "0"
				}
				if evt.Message.GetConversation() != "" {
					message := fmt.Sprintf("%s", evt.Message.GetConversation())
					json_message, _ := json.Marshal(message)
					json_data := fmt.Sprintf(`{"type":"message","sender_jid":"%s", "receiver_jid":"%s", "sender_pushname":%s, "is_from_myself":"%s", "is_group":"%s", "message":%s, "time_stamp":%s, "message_id":"%s"}`, sender_jid, receiver_jid, json_sender_pushname, is_from_myself, is_group, json_message, json_time_stamp, message_id)
					log.Infof("%s", json_data)
					if tasker_type == "net.dinglisch.android.taskerm" {
						//intentTaskerm := exec.Command("/system/bin/sh", "/data/data/net.dinglisch.android.taskerm/files/whatsmeow3/tasker-am/com.taskerm.termuxam.am", "broadcast", "-a", "intent.from.mdtest.v3", "-p", "net.dinglisch.android.taskerm", "--es", "json_data", json_data)
						//intentTaskerm.Output()
						file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.taskerm/files/whatsmeow3/event", "message_")
						encoded_json_data := base64.StdEncoding.EncodeToString([]byte(json_data))
						_, _ = file.WriteString(encoded_json_data)
						_ = file.Close()
					} else if tasker_type == "net.dinglisch.android.tasker" {
						//intentTasker := exec.Command("/system/bin/sh", "/data/data/net.dinglisch.android.tasker/files/whatsmeow3/tasker-am/com.tasker.termuxam.am", "broadcast", "-a", "intent.from.mdtest.v3", "-p", "net.dinglisch.android.tasker", "--es", "json_data", json_data)
						//intentTasker.Output()
						file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.tasker/files/whatsmeow3/event", "message_")
						encoded_json_data := base64.StdEncoding.EncodeToString([]byte(json_data))
						_, _ = file.WriteString(encoded_json_data)
						_ = file.Close()
					} else {
						intentTaskerm := exec.Command("am", "broadcast", "-a", "intent.from.mdtest.v3", "-p", "net.dinglisch.android.taskerm", "--es", "json_data", json_data)
						intentTaskerm.Output()
						intentTasker := exec.Command("am", "broadcast", "-a", "intent.from.mdtest.v3", "-p", "net.dinglisch.android.tasker", "--es", "json_data", json_data)
						intentTasker.Output()
					}
				} else if evt.Message.GetExtendedTextMessage() != nil {
					extended_message := fmt.Sprintf("%s", evt.Message.ExtendedTextMessage.GetText())
					json_extended_message, _ := json.Marshal(extended_message)
					json_data := fmt.Sprintf(`{"type":"extended_message","sender_jid":"%s", "receiver_jid":"%s", "sender_pushname":%s, "is_from_myself":"%s", "is_group":"%s", "message":%s, "time_stamp":%s, "message_id":"%s"}`, sender_jid, receiver_jid, json_sender_pushname, is_from_myself, is_group, json_extended_message, json_time_stamp, message_id)
					log.Infof("%s", json_data)
					if tasker_type == "net.dinglisch.android.taskerm" {
						file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.taskerm/files/whatsmeow3/event", "message_")
						encoded_json_data := base64.StdEncoding.EncodeToString([]byte(json_data))
						_, _ = file.WriteString(encoded_json_data)
						_ = file.Close()
					} else if tasker_type == "net.dinglisch.android.tasker" {
						file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.tasker/files/whatsmeow3/event", "message_")
						encoded_json_data := base64.StdEncoding.EncodeToString([]byte(json_data))
						_, _ = file.WriteString(encoded_json_data)
						_ = file.Close()
					} else {
						intentTaskerm := exec.Command("am", "broadcast", "-a", "intent.from.mdtest.v3", "-p", "net.dinglisch.android.taskerm", "--es", "json_data", json_data)
						intentTaskerm.Output()
						intentTasker := exec.Command("am", "broadcast", "-a", "intent.from.mdtest.v3", "-p", "net.dinglisch.android.tasker", "--es", "json_data", json_data)
						intentTasker.Output()
					}
				} else if evt.Message.GetButtonsResponseMessage() != nil {
					origin_message_id := fmt.Sprintf("%s", evt.Message.ButtonsResponseMessage.ContextInfo.GetStanzaId())
					button_selected_button := fmt.Sprintf("%s", evt.Message.ButtonsResponseMessage.GetSelectedDisplayText())
					json_button_selected_button, _ := json.Marshal(button_selected_button)
					button_title := fmt.Sprintf("%s", evt.Message.ButtonsResponseMessage.ContextInfo.QuotedMessage.ButtonsMessage.GetText())
					json_button_title, _ := json.Marshal(button_title)
					button_body := fmt.Sprintf("%s", evt.Message.ButtonsResponseMessage.ContextInfo.QuotedMessage.ButtonsMessage.GetContentText())
					json_button_body, _ := json.Marshal(button_body)
					button_footer := fmt.Sprintf("%s", evt.Message.ButtonsResponseMessage.ContextInfo.QuotedMessage.ButtonsMessage.GetFooterText())
					json_button_footer, _ := json.Marshal(button_footer)
					json_data := fmt.Sprintf(`{"type":"button_response_message","sender_jid":"%s", "receiver_jid":"%s", "sender_pushname":%s, "is_from_myself":"%s", "is_group":"%s", "button_selected_button":%s, "button_title":%s , "button_body":%s, "button_footer":%s, "time_stamp":%s, "message_id":"%s", "origin_message_id":"%s"}`, sender_jid, receiver_jid, json_sender_pushname, is_from_myself, is_group, json_button_selected_button, json_button_title, json_button_body, json_button_footer, json_time_stamp, message_id, origin_message_id)
					log.Infof("%s", json_data)
					if tasker_type == "net.dinglisch.android.taskerm" {
						file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.taskerm/files/whatsmeow3/event", "message_")
						encoded_json_data := base64.StdEncoding.EncodeToString([]byte(json_data))
						_, _ = file.WriteString(encoded_json_data)
						_ = file.Close()
					} else if tasker_type == "net.dinglisch.android.tasker" {
						file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.tasker/files/whatsmeow3/event", "message_")
						encoded_json_data := base64.StdEncoding.EncodeToString([]byte(json_data))
						_, _ = file.WriteString(encoded_json_data)
						_ = file.Close()
					} else {
						intentTaskerm := exec.Command("am", "broadcast", "-a", "intent.from.mdtest.v3", "-p", "net.dinglisch.android.taskerm", "--es", "json_data", json_data)
						intentTaskerm.Output()
						intentTasker := exec.Command("am", "broadcast", "-a", "intent.from.mdtest.v3", "-p", "net.dinglisch.android.tasker", "--es", "json_data", json_data)
						intentTasker.Output()
					}
				} else if evt.Message.GetListResponseMessage() != nil {
					
					origin_message_id := fmt.Sprintf("%s", evt.Message.ListResponseMessage.ContextInfo.GetStanzaId())
					
					list_selected_title := fmt.Sprintf("%s", evt.Message.ListResponseMessage.GetTitle())
					json_list_selected_title, _ := json.Marshal(list_selected_title)
					
					list_selected_description := fmt.Sprintf("%s", evt.Message.ListResponseMessage.GetDescription())
					json_list_selected_description, _ := json.Marshal(list_selected_description)
					
					list_title := fmt.Sprintf("%s", evt.Message.ListResponseMessage.ContextInfo.QuotedMessage.ListMessage.GetTitle())
					json_list_title, _ := json.Marshal(list_title)
					
					list_body := fmt.Sprintf("%s", evt.Message.ListResponseMessage.ContextInfo.QuotedMessage.ListMessage.GetDescription())
					json_list_body, _ := json.Marshal(list_body)
					
					list_footer := fmt.Sprintf("%s", evt.Message.ListResponseMessage.ContextInfo.QuotedMessage.ListMessage.GetFooterText())
					json_list_footer, _ := json.Marshal(list_footer)
					
					list_button_text := fmt.Sprintf("%s", evt.Message.ListResponseMessage.ContextInfo.QuotedMessage.ListMessage.GetButtonText())
					json_list_button_text, _ := json.Marshal(list_button_text)
					
					list_header := fmt.Sprintf("%s", evt.Message.ListResponseMessage.ContextInfo.QuotedMessage.ListMessage.Sections[0].GetTitle())
					json_list_header, _ := json.Marshal(list_header)
					
					json_data := fmt.Sprintf(`{"type":"list_response_message","sender_jid":"%s", "receiver_jid":"%s", "sender_pushname":%s, "is_from_myself":"%s", "is_group":"%s", "list_selected_title":%s, "list_selected_description":%s, "list_title":%s, "list_body":%s, "list_footer":%s, "list_button_text":%s, "list_header":%s, "time_stamp":%s, "origin_message_id":"%s", "message_id":"%s"}`, sender_jid, receiver_jid, json_sender_pushname, is_from_myself, is_group, json_list_selected_title, json_list_selected_description, json_list_title, json_list_body, json_list_footer, json_list_button_text, json_list_header, json_time_stamp, origin_message_id, message_id)
					log.Infof("%s", json_data)
					if tasker_type == "net.dinglisch.android.taskerm" {
						file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.taskerm/files/whatsmeow3/event", "message_")
						encoded_json_data := base64.StdEncoding.EncodeToString([]byte(json_data))
						_, _ = file.WriteString(encoded_json_data)
						_ = file.Close()
					} else if tasker_type == "net.dinglisch.android.tasker" {
						file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.tasker/files/whatsmeow3/event", "message_")
						encoded_json_data := base64.StdEncoding.EncodeToString([]byte(json_data))
						_, _ = file.WriteString(encoded_json_data)
						_ = file.Close()
					} else {
						intentTaskerm := exec.Command("am", "broadcast", "-a", "intent.from.mdtest.v3", "-p", "net.dinglisch.android.taskerm", "--es", "json_data", json_data)
						intentTaskerm.Output()
						intentTasker := exec.Command("am", "broadcast", "-a", "intent.from.mdtest.v3", "-p", "net.dinglisch.android.tasker", "--es", "json_data", json_data)
						intentTasker.Output()
					}
				} else if evt.Message.GetPollUpdateMessage() != nil {
					message_id = fmt.Sprintf("%s", evt.Message.PollUpdateMessage.PollCreationMessageKey.GetId())
					decrypted, err := cli.DecryptPollVote(evt)
					if err != nil {
						//log.Errorf("Failed to decrypt vote: %v", err)
					} else {
						selected_options_arr := []string{}
						for _, option := range decrypted.SelectedOptions {
							selected_option := fmt.Sprintf("%X", option)
							selected_options_arr = append(selected_options_arr, selected_option)
						}
						selected_options := strings.Join(selected_options_arr, ",")
						json_selected_options, _ := json.Marshal(selected_options)
						//json_data := fmt.Sprintf(`{"type":"poll_update_message","sender_jid":"%s", "receiver_jid":"%s", "sender_pushname":%s, "is_from_myself":"%s", "is_group":"%s", "poll_selected_options":%s, "poll_id":"%s", "time_stamp":%s, "message_id":"%s"}`, sender_jid, receiver_jid, json_sender_pushname, is_from_myself, is_group, json_selected_options, evt.Message.PollUpdateMessage.PollCreationMessageKey.GetId(), json_time_stamp, message_id)
						json_data := fmt.Sprintf(`{"type":"poll_response_message","sender_jid":"%s", "receiver_jid":"%s", "sender_pushname":%s, "is_from_myself":"%s", "is_group":"%s", "poll_selected_options":%s, "time_stamp":%s, "message_id":"%s"}`, sender_jid, receiver_jid, json_sender_pushname, is_from_myself, is_group, json_selected_options, json_time_stamp, message_id)
						log.Infof("%s", json_data)
						//log.Infof("%s", evt.Message.PollUpdateMessage.PollCreationMessageKey.GetId())
						if tasker_type == "net.dinglisch.android.taskerm" {
							file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.taskerm/files/whatsmeow3/event", "message_")
							encoded_json_data := base64.StdEncoding.EncodeToString([]byte(json_data))
							_, _ = file.WriteString(encoded_json_data)
							_ = file.Close()
						} else if tasker_type == "net.dinglisch.android.tasker" {
							file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.tasker/files/whatsmeow3/event", "message_")
							encoded_json_data := base64.StdEncoding.EncodeToString([]byte(json_data))
							_, _ = file.WriteString(encoded_json_data)
							_ = file.Close()
						} else {
							intentTaskerm := exec.Command("am", "broadcast", "-a", "intent.from.mdtest.v3", "-p", "net.dinglisch.android.taskerm", "--es", "json_data", json_data)
							intentTaskerm.Output()
							intentTasker := exec.Command("am", "broadcast", "-a", "intent.from.mdtest.v3", "-p", "net.dinglisch.android.tasker", "--es", "json_data", json_data)
							intentTasker.Output()
						}
					}
				} else if evt.Message.GetPollCreationMessage() != nil || evt.Message.GetPollCreationMessageV2() != nil || evt.Message.GetPollCreationMessageV3() != nil {
					var question string
					type Option struct {
						Option_Name string
					}
					var option Option
					if evt.Message.GetPollCreationMessage() != nil  {
						question = fmt.Sprintf("%s", evt.Message.PollCreationMessage.GetName())
						for _, optionName := range evt.Message.PollCreationMessage.GetOptions() {
							option_json := `{"Option_Name"` + strings.TrimPrefix(fmt.Sprintf("%s", optionName), "optionName") + `}`
							json.Unmarshal([]byte(option_json), &option)
							//option := strings.TrimPrefix(fmt.Sprintf("%s", optionName), "optionName:")
							sum := sha256.Sum256([]byte(option.Option_Name))
							//sum := sha256.Sum256([]byte(strings.TrimSuffix(strings.TrimPrefix(fmt.Sprintf("%s", optionName), `optionName:"`), `"`)))
							sha := fmt.Sprintf("%x", sum)
							log.Infof("%s = %s", option.Option_Name, sha)
							//log.Infof("%s = %s", option, sha)
							if tasker_type == "net.dinglisch.android.taskerm" {
								file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.taskerm/files/whatsmeow3/.tmp", sha + "_poll_option_")
								encoded_json_data := base64.StdEncoding.EncodeToString([]byte(option.Option_Name))
								_, _ = file.WriteString(encoded_json_data)
								_ = file.Close()
							} else if tasker_type == "net.dinglisch.android.tasker" {
								file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.tasker/files/whatsmeow3/.tmp", sha + "_poll_option_")
								encoded_json_data := base64.StdEncoding.EncodeToString([]byte(option.Option_Name))
								_, _ = file.WriteString(encoded_json_data)
								_ = file.Close()
							}
						}
					} else if evt.Message.GetPollCreationMessageV2() != nil {
						question = fmt.Sprintf("%s", evt.Message.PollCreationMessageV2.GetName())
						for _, optionName := range evt.Message.PollCreationMessageV2.GetOptions() {
							option_json := `{"Option_Name"` + strings.TrimPrefix(fmt.Sprintf("%s", optionName), "optionName") + `}`
							json.Unmarshal([]byte(option_json), &option)
							sum := sha256.Sum256([]byte(option.Option_Name))
							sha := fmt.Sprintf("%x", sum)
							log.Infof("%s = %s", option.Option_Name, sha)
							if tasker_type == "net.dinglisch.android.taskerm" {
								file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.taskerm/files/whatsmeow3/.tmp", sha + "_poll_option_")
								encoded_json_data := base64.StdEncoding.EncodeToString([]byte(option.Option_Name))
								_, _ = file.WriteString(encoded_json_data)
								_ = file.Close()
							} else if tasker_type == "net.dinglisch.android.tasker" {
								file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.tasker/files/whatsmeow3/.tmp", sha + "_poll_option_")
								encoded_json_data := base64.StdEncoding.EncodeToString([]byte(option.Option_Name))
								_, _ = file.WriteString(encoded_json_data)
								_ = file.Close()
							}
						}
					} else if evt.Message.GetPollCreationMessageV3() != nil {
						question = fmt.Sprintf("%s", evt.Message.PollCreationMessageV3.GetName())
						for _, optionName := range evt.Message.PollCreationMessageV3.GetOptions() {
							option_json := `{"Option_Name"` + strings.TrimPrefix(fmt.Sprintf("%s", optionName), "optionName") + `}`
							json.Unmarshal([]byte(option_json), &option)
							sum := sha256.Sum256([]byte(option.Option_Name))
							sha := fmt.Sprintf("%x", sum)
							log.Infof("%s = %s", option.Option_Name, sha)
							if tasker_type == "net.dinglisch.android.taskerm" {
								file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.taskerm/files/whatsmeow3/.tmp", sha + "_poll_option_")
								encoded_json_data := base64.StdEncoding.EncodeToString([]byte(option.Option_Name))
								_, _ = file.WriteString(encoded_json_data)
								_ = file.Close()
							} else if tasker_type == "net.dinglisch.android.tasker" {
								file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.tasker/files/whatsmeow3/.tmp", sha + "_poll_option_")
								encoded_json_data := base64.StdEncoding.EncodeToString([]byte(option.Option_Name))
								_, _ = file.WriteString(encoded_json_data)
								_ = file.Close()
							}
						}
					}
					
					json_question, _ := json.Marshal(question)
					if tasker_type == "net.dinglisch.android.taskerm" {
						file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.taskerm/files/whatsmeow3/.tmp", message_id + "_" + receiver_jid + "_poll_question_")
						encoded_json_data := base64.StdEncoding.EncodeToString([]byte(json_question))
						_, _ = file.WriteString(encoded_json_data)
						_ = file.Close()
					} else if tasker_type == "net.dinglisch.android.tasker" {
						file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.taskerm/files/whatsmeow3/.tmp", message_id + "_" + receiver_jid + "_poll_question_")
						encoded_json_data := base64.StdEncoding.EncodeToString([]byte(json_question))
						_, _ = file.WriteString(encoded_json_data)
						_ = file.Close()
					}
				} else if (evt.Message.GetPollCreationMessage() != nil || evt.Message.GetPollCreationMessageV2() != nil || evt.Message.GetPollCreationMessageV3() != nil) && evt.Message.GetConversation() == "null" {
					var question string
					var multiple_choice string
					var selectable_options_count string
					options_arr := []string{}
					if evt.Message.GetPollCreationMessage() != nil  {
						question = fmt.Sprintf("%s", evt.Message.PollCreationMessage.GetName())
						selectable_options_count = fmt.Sprintf("%d", evt.Message.PollCreationMessage.GetSelectableOptionsCount())
						for _, optionName := range evt.Message.PollCreationMessage.GetOptions() {
							option := fmt.Sprintf("%s", optionName)
							options_arr = append(options_arr, option)
						}
					} else if evt.Message.GetPollCreationMessageV2() != nil {
						question = fmt.Sprintf("%s", evt.Message.PollCreationMessageV2.GetName())
						selectable_options_count = fmt.Sprintf("%d", evt.Message.PollCreationMessageV2.GetSelectableOptionsCount())
						for _, optionName := range evt.Message.PollCreationMessageV2.GetOptions() {
							option := fmt.Sprintf("%s", optionName)
							options_arr = append(options_arr, option)
						}
					} else if evt.Message.GetPollCreationMessageV3() != nil {
						question = fmt.Sprintf("%s", evt.Message.PollCreationMessageV2.GetName())
						selectable_options_count = fmt.Sprintf("%d", evt.Message.PollCreationMessageV2.GetSelectableOptionsCount())
						for _, optionName := range evt.Message.PollCreationMessageV2.GetOptions() {
							option := fmt.Sprintf("%s", optionName)
							options_arr = append(options_arr, option)
						}
					}
					if selectable_options_count == "1" {
						multiple_choice = "0"
					} else {
						multiple_choice = "1"
					}
					json_question, _ := json.Marshal(question)
					json_data := fmt.Sprintf(`{"type":"poll_message","sender_jid":"%s", "receiver_jid":"%s", "sender_pushname":%s, "is_from_myself":"%s", "is_group":"%s", "question":%s, "options":[{%s}], "multiple_choice":"%s", "time_stamp":%s, "message_id":"%s"}`, sender_jid, receiver_jid, json_sender_pushname, is_from_myself, is_group, json_question, strings.Join(options_arr, "},{"), multiple_choice, json_time_stamp, message_id)
					log.Infof("%s", json_data)
					if tasker_type == "net.dinglisch.android.taskerm" {
						file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.taskerm/files/whatsmeow3/event", "message_")
						encoded_json_data := base64.StdEncoding.EncodeToString([]byte(json_data))
						_, _ = file.WriteString(encoded_json_data)
						_ = file.Close()
					} else if tasker_type == "net.dinglisch.android.tasker" {
						file, _ := ioutil.TempFile("/data/data/net.dinglisch.android.tasker/files/whatsmeow3/event", "message_")
						encoded_json_data := base64.StdEncoding.EncodeToString([]byte(json_data))
						_, _ = file.WriteString(encoded_json_data)
						_ = file.Close()
					} else {
						intentTaskerm := exec.Command("am", "broadcast", "-a", "intent.from.mdtest.v3", "-p", "net.dinglisch.android.taskerm", "--es", "json_data", json_data)
						intentTaskerm.Output()
						intentTasker := exec.Command("am", "broadcast", "-a", "intent.from.mdtest.v3", "-p", "net.dinglisch.android.tasker", "--es", "json_data", json_data)
						intentTasker.Output()
					}
				}
			}
		}()

		if evt.Message.GetPollUpdateMessage() != nil {
			decrypted, err := cli.DecryptPollVote(evt)
			if err != nil {
				log.Errorf("Failed to decrypt vote: %v", err)
			} else {
				log.Infof("Selected options in decrypted vote:")
				for _, option := range decrypted.SelectedOptions {
					log.Infof("- %X", option)
				}
			}
		} else if evt.Message.GetEncReactionMessage() != nil {
			decrypted, err := cli.DecryptReaction(evt)
			if err != nil {
				log.Errorf("Failed to decrypt encrypted reaction: %v", err)
			} else {
				log.Infof("Decrypted reaction: %+v", decrypted)
			}
		}

	case *events.Receipt:
		if evt.Type == events.ReceiptTypeRead || evt.Type == events.ReceiptTypeReadSelf {
			log.Infof("%v was read by %s at %s", evt.MessageIDs, evt.SourceString(), evt.Timestamp)
		} else if evt.Type == events.ReceiptTypeDelivered {
			log.Infof("%s was delivered to %s at %s", evt.MessageIDs[0], evt.SourceString(), evt.Timestamp)
		}
	case *events.Presence:
		if evt.Unavailable {
			if evt.LastSeen.IsZero() {
				log.Infof("%s is now offline", evt.From)
			} else {
				log.Infof("%s is now offline (last seen: %s)", evt.From, evt.LastSeen)
			}
		} else {
			log.Infof("%s is now online", evt.From)
		}
	case *events.HistorySync:
		id := atomic.AddInt32(&historySyncID, 1)
		fileName := fmt.Sprintf("history-%d-%d.json", startupTime, id)
		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Errorf("Failed to open file to write history sync: %v", err)
			return
		}
		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		err = enc.Encode(evt.Data)
		if err != nil {
			log.Errorf("Failed to write history sync: %v", err)
			return
		}
		log.Infof("Wrote history sync to %s", fileName)
		_ = file.Close()
	case *events.AppState:
		log.Debugf("App state event: %+v / %+v", evt.Index, evt.SyncActionValue)
	case *events.KeepAliveTimeout:
		log.Debugf("Keepalive timeout event: %+v", evt)
		if evt.ErrorCount > 3 {
			log.Debugf("Got >3 keepalive timeouts, forcing reconnect")
			go func() {
if is_mode == "both" || is_mode == "receive" || is_mode == "send" {
	kill_server()
}
				cli.Disconnect()
				err := cli.Connect()
				if err != nil {
					log.Errorf("Error force-reconnecting after keepalive timeouts: %v", err)
				}
			}()
		}
	case *events.KeepAliveRestored:
		log.Debugf("Keepalive restored")
	}
}
