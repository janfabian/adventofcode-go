package aoc2024

import (
	"adventofcode/lib"
	"fmt"
	"strconv"
	"strings"
)

type Segment struct {
	Length int
	Pos    int
	Id     int
	Free   bool
}

func CreateSegment(input []int) ([]*Segment, int) {
	segments := []*Segment{}
	isFile := true
	fileId := 0
	pos := 0

	for _, val := range input {
		var s Segment
		if isFile {
			s = Segment{Length: val, Pos: pos, Id: fileId, Free: false}
			fileId++
		} else {
			s = Segment{Length: val, Pos: pos, Id: 0, Free: true}
		}

		pos += val
		isFile = !isFile
		if s.Free && s.Length == 0 {
			continue
		}
		segments = append(segments, &s)
	}

	return segments, pos
}

func Solve(input []int, part2 bool) int {
	if part2 {
		return Part2(input)

	}

	return Part1(input)
}

func Part1(input []int) int {
	segments, pos := CreateSegment(input)
	back_pos := pos
	back_ix := len(segments) - 1
	back_s := segments[back_ix]
	front_ix := 0
	front_s := segments[front_ix]
	result_str := ""
	result_ix := 0
	result := 0

	for front_pos := 0; front_pos < pos; front_pos++ {
		if back_pos <= front_pos {
			break
		}

		if front_s.Pos+front_s.Length-1 < front_pos {
			front_ix++
		}

		front_s = segments[front_ix]

		if !front_s.Free {
			result_str += strconv.Itoa(front_s.Id)
			result += front_s.Id * result_ix
			result_ix++
			continue
		}

		for {
			if back_pos <= front_pos {
				break
			}

			if back_s.Free {
				back_pos -= back_s.Length
				back_ix--
				back_s = segments[back_ix]
			} else {
				result_str += strconv.Itoa(back_s.Id)
				result += back_s.Id * result_ix
				result_ix++
				back_pos--

				if back_s.Pos >= back_pos {
					back_ix--
					back_s = segments[back_ix]
				}

				break
			}
		}

	}

	return result
}

type DoubleLinkedSegment struct {
	S    *Segment
	Next *DoubleLinkedSegment
	Prev *DoubleLinkedSegment
}

func Part2(input []int) int {
	segments, _ := CreateSegment(input)

	var prev *DoubleLinkedSegment
	var head *DoubleLinkedSegment
	for _, s := range segments {

		dls := &DoubleLinkedSegment{S: s, Prev: prev}

		if prev != nil {
			prev.Next = dls
		}

		if head == nil {
			head = dls
		}

		prev = dls

	}

	file_tail := prev

	for file_tail != nil {
		if file_tail.S.Free {
			file_tail = file_tail.Prev
			continue
		}

		free_head := head

		for free_head != file_tail {

			if !free_head.S.Free {
				free_head = free_head.Next
				continue
			}

			if free_head.S.Length >= file_tail.S.Length {
				subst := &DoubleLinkedSegment{
					S: &Segment{
						Length: file_tail.S.Length,
						Pos:    file_tail.S.Pos,
						Id:     file_tail.S.Id,
						Free:   false,
					}, Prev: free_head.Prev, Next: free_head,
				}

				if free_head.Prev != nil {
					free_head.Prev.Next = subst
				}

				free_head.Prev = subst

				free_head.S.Pos += file_tail.S.Length
				free_head.S.Length -= file_tail.S.Length

				file_tail.S.Free = true
				file_tail.S.Id = 0

				break
			}

			free_head = free_head.Next

		}

		file_tail = file_tail.Prev
	}

	result_ix := 0
	result := 0

	for dls := head; dls != nil; dls = dls.Next {
		s := dls.S
		for i := 0; i < s.Length; i++ {
			if !s.Free {
				result += s.Id * result_ix
			}
			result_ix++
		}
	}

	return result
}

func ParseInput(input string, part2 bool) ([]int, error) {
	_, scanner, err := lib.ScanFile(input)

	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	result := []int{}

	scanner.Scan()
	line := scanner.Text()
	parts := strings.Split(line, "")

	for _, part := range parts {
		intval, err := strconv.Atoi(part)

		if err != nil {
			return nil, fmt.Errorf("error converting to int: %v", err)
		}

		result = append(result, intval)
	}

	return result, nil
}
