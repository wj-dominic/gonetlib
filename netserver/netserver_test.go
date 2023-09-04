package netserver

import (
	"testing"
	"time"
)

/*
1. 서버 옵션 설정
- 엔드 포인트
- 최대 동접자 수
- 서버 코어
ㄴTask 기반
ㄴ명령어 기반

2. 서버 시작


*/

func TestConnect(t *testing.T) {
	go Run("127.0.0.1:8888")
	time.Sleep(time.Second * 30)



	

}
