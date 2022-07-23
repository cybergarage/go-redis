// Copyright (C) 2022 Satoshi Konno All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protocol

const (
	errorEmptyMessage            = "message is short (%d)"
	errorUnknownMessageType      = "unknown message type (%c)"
	errorInvalidMessageType      = "invalid message type (%c)"
	errorInvalidMessage          = "invalid message (%s)"
	errorInvalidBulkStringLength = "invalid bulk string length (%d != %d)"
)
