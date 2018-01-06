package utils

/*
		//使用方法：
		var params []interface{}
		for i:=0; i< 5; i++ {
			params = append(params, fmt.Sprintf("params %d", i))
		}

		utils.QueueTask(2, params, func(idx int, param interface{}) {
			utils.JJKPrintln(param)
			time.Sleep(time.Second * 2)
		})
 */
type QueueTaskFunc func(idx int, param interface{})
func QueueTask(maxConcurrenceCount int, params []interface{}, taskFunc QueueTaskFunc)  {
	var totalTaskCount int = len(params)
	var concurrenceMax int = maxConcurrenceCount
	var currentCount int = 0
	var finishedCount int = 0
	curCntChan := make(chan int)
	for {
		//add task
		if currentCount < concurrenceMax {
			startIdx := finishedCount+currentCount
			endIdx := finishedCount + (concurrenceMax-currentCount)
			for i := startIdx; i< endIdx ; i++ {
				if i < totalTaskCount {
					currentCount += 1
					thisParam := params[i]
					//执行任务
					go func(key int, param interface{}) {
						taskFunc(key, param)
						curCntChan <- key
					}(finishedCount, thisParam)
				}
			}
		}

		//waiting task
		select {
		case <-curCntChan:
			finishedCount+=1
			currentCount -= 1
		}

		//end this group task
		if finishedCount == totalTaskCount {
			break
		}
	}
}

