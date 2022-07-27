package dp

import "testing"

func TestDP(t *testing.T) {
	n := way1(8, 2, 4, 4)
	t.Log(n)
}

/**
N: 总共有多少步
start: 从哪一步开始
aim: 目标
k: 走多少步
*/
func way1(N int, start int, aim int, k int) int {
	return process1(start, k, aim, N)
}

/**
cur: 机器人当前来到的位置
rest: 机器人还有 rest 步需要走
aim: 最终的目标
N: 有哪些位置
return: 机器人从 cur 出发，走过 rest 步之后，最终停在 aim的方法，是多少。
*/
func process1(cur, rest, aim, N int) int {
	if rest == 0 {
		if cur == aim {
			return 1
		}
		return 0
	}
	if cur == 1 {
		return process1(2, rest-1, aim, N)
	}
	if cur == N {
		return process1(N-1, rest-1, aim, N)
	}
	return process1(cur-1, rest-1, aim, N) + process1(cur+1, rest-1, aim, N)
}
