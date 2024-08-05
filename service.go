package main

import (
	"bufio"
	"errors"
	"fmt"
	"reflect"
	"strings"

	//"fmt"
	"io"
	"log"
)

type JSON map[string]interface{}

type Unknown interface{}

type TRACKER map[int]interface{}

func findCurrent(depth int, tracker TRACKER) int {
	current := depth - 1
	for tracker[current] == nil && current > -1 {
		current -= 1
	}
	//fmt.Println("reading current and tracker ", tracker, current, depth)
	return current
}

func BuildJSON(r *bufio.Reader) JSON {
	tracker := make(TRACKER)

	builder := JSON{}

	tracker[-1] = builder

	ReadYAML(r, tracker)

	return builder
}

func ReadYAML(reader *bufio.Reader, tracker TRACKER) *JSON {
	//
	var bts string
	var key_name string
	var last_record string
	var last_key string
	var depth int
	var is_array bool
	var new_array_start bool
	var array_key string
	var array_keys []string
	var is_value_now bool
	//var previous string

	for {
		b, err := reader.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("losing games ", key_name, last_record, depth)
				if key_name != "" || last_record != "" {
					current := findCurrent(depth, tracker)
					parent_p := findCurrent(current, tracker)

					// fmt.Println("after end of file ", current, parent_p, tracker[current])

					if tracker[parent_p] == nil {
						return nil
					}
					parent_is_array := reflect.TypeOf(tracker[parent_p]).Kind() == reflect.Slice

					fmt.Println("current ", tracker[current])
					if reflect.TypeOf(tracker[current]).Kind() != reflect.Slice {
						tracker[current].(JSON)[strings.TrimSpace(last_record)] = key_name
					} else {
						if len(tracker[current].([]JSON)) > 0 {
							last := tracker[current].([]JSON)[len(tracker[current].([]JSON))-1]
							last[strings.TrimSpace(last_record)] = key_name
						} else {
							key := strings.TrimSpace(last_record)

							if key == "" {
								key = key_name
							}
							if parent_is_array {
								tracker[parent_p].([]JSON)[len(tracker[parent_p].([]JSON))-1][array_key] = append(
									tracker[parent_p].([]JSON)[len(tracker[parent_p].([]JSON))-1][array_key].([]JSON), JSON{
										strings.TrimSpace(key): strings.TrimSpace(key_name),
									})
							} else {
								tracker[current] = append(tracker[current].([]JSON), JSON{
									last_record: key_name,
								})
							}
						}
					}
				}
				break
			}
			log.Fatal("Failed to read YAML doc")
		}
		bts = string(b)

		if bts == "\n" {
			// check the depth
			current := findCurrent(depth, tracker)
			if current < -1 {
				log.Fatal("trouble")
			}
			//ldepth = depth
			parent := tracker[current]
			is_array = reflect.TypeOf(parent).Kind() == reflect.Slice
			parent_p := findCurrent(current, tracker)

			var is_parent_parent_array bool = false

			if tracker[parent_p] != nil {
				is_parent_parent_array = reflect.TypeOf(tracker[parent_p]).Kind() == reflect.Slice
			}

			key_name = strings.TrimSpace(key_name)
			if len(key_name) > 0 {
				if is_array {
					if is_parent_parent_array {
						cur_obj := tracker[parent_p].([]JSON)[len(tracker[parent_p].([]JSON))-1][array_key]
						if len(cur_obj.([]JSON)) == 0 || new_array_start {
							//fmt.Println("asda ", cur_obj)
							var key string = last_record
							if !is_value_now {
								key = key_name
							}
							tracker[parent_p].([]JSON)[len(tracker[parent_p].([]JSON))-1][array_key] = append(
								tracker[parent_p].([]JSON)[len(tracker[parent_p].([]JSON))-1][array_key].([]JSON),
								JSON{
									strings.TrimSpace(key): key_name,
								},
							)
							tracker[current] = tracker[parent_p].([]JSON)[len(tracker[parent_p].([]JSON))-1][array_key]

						} else {
							cur := cur_obj.([]JSON)[len(cur_obj.([]JSON))-1]
							var key string = last_record
							if !is_value_now {
								key = key_name
							}
							cur[strings.TrimSpace(key)] = key_name
							tracker[current] = cur_obj
							//fmt.Println("nothing ", cur)
						}

					} else {

						arr_key_len := len(array_keys) - 1
						for i := arr_key_len; i >= 0; i -= 1 {
							array_key = array_keys[i]
							if tracker[parent_p].(JSON)[array_key] != nil {
								break
							}
						}
						// for _, v := range array_keys {
						// 	array_key = v
						// 	if tracker[parent_p].(JSON)[array_key] != nil {
						// 		break
						// 	}
						// }
						if strings.TrimSpace(key_name) == "shared-data" {
							fmt.Println("dasd ", depth, is_array, tracker[parent_p].(JSON)[array_key])
						}
						if len(tracker[parent_p].(JSON)[array_key].([]JSON)) == 0 || new_array_start {
							tracker[parent_p].(JSON)[array_key] = append(tracker[parent_p].(JSON)[array_key].([]JSON), JSON{
								strings.TrimSpace(last_record): key_name,
							})
						} else {
							cur := tracker[parent_p].(JSON)[array_key].([]JSON)[len(tracker[parent_p].(JSON)[array_key].([]JSON))-1]
							cur[strings.TrimSpace(last_record)] = key_name
						}
					}
					if !is_parent_parent_array {
						tracker[current] = tracker[parent_p].(JSON)[array_key]
					}
				} else {
					var key string = last_record
					if !is_value_now {
						key = key_name
					}
					parent.(JSON)[strings.TrimSpace(key)] = key_name
				}
				//fmt.Println("parent ", parent)
			} else if last_record != "" {
				last_key = strings.TrimSpace(last_record)
				new_map := make(JSON)
				tracker[depth] = new_map
				if is_array {
					//parent_p := findCurrent(current, tracker)
					if len(tracker[current].([]JSON)) > 0 {
						cur := tracker[current].([]JSON)[len(tracker[current].([]JSON))-1]
						if is_value_now {
							cur[strings.TrimSpace(last_record)] = tracker[depth]
						} else {
							cur[strings.TrimSpace(last_record)] = strings.TrimSpace(last_record)
						}
					} else {
						if is_value_now {
							tracker[current] = append(tracker[current].([]JSON), JSON{
								strings.TrimSpace(last_record): tracker[depth],
							})
						} else {
							tracker[current] = append(tracker[current].([]JSON), JSON{
								strings.TrimSpace(key_name): strings.TrimSpace(last_record),
							})
						}
						//fmt.Println("whast asdad ", tracker[parent_p])
						tracker[parent_p].(JSON)[array_key] = tracker[current]
					}
				} else {
					parent.(JSON)[strings.TrimSpace(last_record)] = tracker[depth]
				}
			}
			depth = 0
			key_name = ""
			last_record = ""
			is_value_now = false
			new_array_start = false
			// need to find the parent
		} else if bts == ":" && !is_value_now {
			// time to save the key and start getting value
			last_record = key_name
			key_name = ""
			is_value_now = true
			continue
		} else if (bts == " " || bts == "-") && !is_value_now {
			depth += 1
			if bts == "-" {
				// we have encountered a possible
				new_array_start = true
				current := findCurrent(depth, tracker)
				parent_cur := findCurrent(current, tracker)
				array_key = last_key
				array_keys = append(array_keys, array_key)
				//fmt.Println("though ", tracker[parent_cur], tracker[current], last_key)
				cur_obj := tracker[parent_cur]
				new_array := make([]JSON, 0)
				if parent_cur >= -1 {
					if reflect.TypeOf(cur_obj).Kind() != reflect.Slice {
						if cur_obj.(JSON)[last_key] != nil && reflect.TypeOf(cur_obj.(JSON)[last_key]).Kind() != reflect.Slice {
							tracker[parent_cur].(JSON)[last_key] = new_array
							tracker[current] = new_array
						}
					} else {
						item := cur_obj.([]JSON)[len(cur_obj.([]JSON))-1]
						if item[array_key] != nil {
							if reflect.TypeOf(item[strings.TrimSpace(array_key)]).Kind() != reflect.Slice {
								item[array_key] = new_array

							}
						} else {
							item[array_key] = new_array
						}
						tracker[current] = item[array_key]
					}
				}
				continue
			}
		}
		key_name += bts
	}
	return nil
}
