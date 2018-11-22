module.exports = {
  "nested": {
    "protocol": {
      "nested": {
        "MessageType": {
          "values": {
            "MSG_TYPE_CHAT_MESSGAE_REQ": 1,
            "MSG_TYPE_CHAT_MESSAGE_RES": 2,
            "MSG_TYPE_SEARCH_A_GAME_REQ": 3,
            "MSG_TYPE_SEARCH_A_GAME_RES": 4,
            "MSG_TYPE_STOP_SEARCH_REQ": 5,
            "MSG_TYPE_STOP_SEARCH_RES": 6,
            "MSG_TYPE_START_RES": 8,
            "MSG_TYPE_LINE_A_POINT_REQ": 9,
            "MSG_TYPE_LINE_A_POINT_RES": 10,
            "MSG_TYPE_LINE_A_POINT_TO_REQUEST_RES": 14,
            "MSG_TYPE_END_GAME_REQ": 11,
            "MSG_TYPE_END_GAME_RES": 12,
            "MSG_TYPE_CREATE_USER_REQ": 101,
            "MSG_TYPE_CREATE_USER_RES": 102,
            "MSG_TYPE_LOGIN_REQ": 103,
            "MSG_TYPE_LOGIN_RES": 104,
            "MSG_TYPE_LOGOUT_REQ": 105,
            "MSG_TYPE_LOGOUT_RES": 106
          }
        },
        "MobileSuiteModel": {
          "fields": {
            "type": {
              "rule": "required",
              "type": "int32",
              "id": 1
            },
            "message": {
              "type": "bytes",
              "id": 2
            }
          }
        },
        "ChatMsg": {
          "fields": {
            "chatType": {
              "rule": "required",
              "type": "int32",
              "id": 1
            },
            "userId": {
              "rule": "required",
              "type": "int64",
              "id": 2
            },
            "uName": {
              "rule": "required",
              "type": "string",
              "id": 3
            },
            "chatContext": {
              "rule": "required",
              "type": "string",
              "id": 4
            }
          }
        },
        "GameStartDTO": {
          "fields": {
            "opptName": {
              "rule": "required",
              "type": "string",
              "id": 1
            },
            "playerIndex": {
              "rule": "required",
              "type": "int32",
              "id": 2
            }
          }
        },
        "LineAPointDTO": {
          "fields": {
            "row": {
              "rule": "required",
              "type": "int32",
              "id": 1
            },
            "col": {
              "rule": "required",
              "type": "int32",
              "id": 2
            },
            "playerIndex": {
              "rule": "required",
              "type": "int32",
              "id": 3
            }
          }
        },
        "LineAPointResponseDTO": {
          "fields": {
            "result": {
              "rule": "required",
              "type": "int32",
              "id": 1
            }
          }
        },
        "CreateUserDTO": {
          "fields": {
            "uName": {
              "rule": "required",
              "type": "string",
              "id": 1
            },
            "pwd": {
              "rule": "required",
              "type": "string",
              "id": 2
            }
          }
        },
        "CreateResultDTO": {
          "fields": {
            "userId": {
              "rule": "required",
              "type": "int64",
              "id": 1
            }
          }
        },
        "LoginDTO": {
          "fields": {
            "userId": {
              "rule": "required",
              "type": "int64",
              "id": 1
            },
            "uName": {
              "rule": "required",
              "type": "string",
              "id": 2
            },
            "pwd": {
              "rule": "required",
              "type": "string",
              "id": 3
            }
          }
        },
        "LoginResultDTO": {
          "fields": {
            "userId": {
              "rule": "required",
              "type": "int64",
              "id": 1
            },
            "uName": {
              "rule": "required",
              "type": "string",
              "id": 2
            },
            "round": {
              "rule": "required",
              "type": "int32",
              "id": 3
            },
            "winCount": {
              "rule": "required",
              "type": "int32",
              "id": 4
            },
            "rank": {
              "rule": "required",
              "type": "int32",
              "id": 5
            }
          }
        },
        "LogoutDTO": {
          "fields": {
            "userId": {
              "rule": "required",
              "type": "int64",
              "id": 1
            }
          }
        }
      }
    }
  }
}