package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/koala-labs/outback/pkg/outback"

	lru "github.com/hashicorp/golang-lru"
	Outback "github.com/koala-labs/outback/pkg/outback"
	"github.com/spf13/cobra"
)

type Empty struct{}

type LogsOperation struct {
	LogGroupName   string
	Namespace      string
	Service        string
	EndTime        time.Time
	StartTime      time.Time
	Filter         string
	Follow         bool
	LogStreamNames []string
	EventCache     *lru.Cache
}

const (
	timeFormat          = "2006-01-02 15:04:05"
	timeFormatWithZone  = "2006-01-02 15:04:05 MST"
	logStreamNameFormat = "%s/%s/%s"
	eventCacheSize      = 10000
)

var (
	flagServiceLogsFilter    string
	flagServiceLogsEndTime   string
	flagServiceLogsStartTime string
	flagServiceLogsFollow    bool
	flagServiceLogsTasks     []string
)

var serviceLogsCmd = &cobra.Command{
	Use:   "logs <service-name>",
	Short: "Show logs from tasks in a service",
	Long: `Show logs from tasks in a service

	Return either a specific segment of service logs or tail logs in real-time
	using the --follow option. Logs are prefixed by their log stream name which is
	in the format of "fargate/\<service-name>/\<task-id>."
	Follow will continue to run and return logs until interrupted by Control-C. If
	--follow is passed --end cannot be specified.
	Logs can be returned for specific tasks within a service by passing a task ID
	via the --task flag. Pass --task with a task ID multiple times in order to
	retrieve logs from multiple specific tasks.
	A specific window of logs can be requested by passing --start and --end options
	with a time expression. The time expression can be either a duration or a
	timestamp:
	  - Duration (e.g. -1h [one hour ago], -1h10m30s [one hour, ten minutes, and
		thirty seconds ago], 2h [two hours from now])
	  - Timestamp with optional timezone in the format of YYYY-MM-DD HH:MM:SS [TZ];
		timezone will default to UTC if omitted (e.g. 2017-12-22 15:10:03 EST)
	You can filter logs for specific term by passing a filter expression via the
	--filter flag. Pass a single term to search for that term, pass multiple terms
	to search for log messages that include all terms.`,
	Args: cobra.ExactArgs(1),
	Run:  getOrFollowLogs,
}

func getOrFollowLogs(cmd *cobra.Command, args []string) {
	o := &LogsOperation{
		LogGroupName: args[0],
		Filter:       flagServiceLogsFilter,
		Follow:       flagServiceLogsFollow,
		Namespace:    args[0],
		Service:      "",
	}

	o.AddTasks(flagServiceLogsTasks)
	o.AddStartTime(flagServiceLogsStartTime)
	o.AddEndTime(flagServiceLogsEndTime)

	if flagServiceLogsFollow {
		followLogs(o)
	} else {
		getLogs(o)
	}
}

func (o *LogsOperation) AddStartTime(rawStartTime string) {
	if rawStartTime != "" {
		o.StartTime, _ = o.parseTime(rawStartTime)
	}
}

func (o *LogsOperation) AddEndTime(rawEndTime string) {
	if rawEndTime != "" {
		o.EndTime, _ = o.parseTime(rawEndTime)
	}
}

func (o *LogsOperation) AddTasks(tasks []string) {
	for _, task := range tasks {
		logStreamName := fmt.Sprintf(logStreamNameFormat, o.Namespace, o.Service, task)
		o.LogStreamNames = append(o.LogStreamNames, logStreamName)
	}
}

func (o *LogsOperation) SeenEvent(eventID string) bool {
	if o.EventCache == nil {
		o.EventCache, _ = lru.New(eventCacheSize)
	}

	if !o.EventCache.Contains(eventID) {
		o.EventCache.Add(eventID, Empty{})
		return false
	}

	return true
}

func (o *LogsOperation) Validate() error {
	if o.Follow && o.EndTime.IsZero() {
		return ErrCantFollowWithEndTime
	}

	return nil
}

func (o *LogsOperation) parseTime(rawTime string) (time.Time, error) {
	var t time.Time

	if duration, err := time.ParseDuration(strings.ToLower(rawTime)); err == nil {
		return time.Now().Add(duration), nil
	}

	if t, err := time.Parse(timeFormat, rawTime); err == nil {
		return t, nil
	}

	if t, err := time.Parse(timeFormatWithZone, rawTime); err == nil {
		return t, nil
	}

	return t, ErrCouldNotParseTime
}

func followLogs(o *LogsOperation) {
	ticker := time.NewTicker(time.Second)

	if o.StartTime.IsZero() {
		o.StartTime = time.Now()
	}

	for {
		getLogs(o)

		if newStartTime := time.Now().Add(-10 * time.Second); newStartTime.After(o.StartTime) {
			o.StartTime = newStartTime
		}

		<-ticker.C
	}
}

func getLogs(o *LogsOperation) {
	u := Outback.New(awsConfig)

	in := &outback.GetLogsInput{
		LogStreamNames: o.LogStreamNames,
		LogGroupName:   o.LogGroupName,
		Filter:         o.Filter,
		StartTime:      o.StartTime,
		EndTime:        o.EndTime,
	}

	logs, _ := u.GetLogs(in)

	for _, logLine := range logs {
		if !o.SeenEvent(logLine.EventID) {
			fmt.Printf("[%s][%s] - %s\n", logLine.Timestamp, logLine.LogStreamName, logLine.Message)
		}
	}
}

func init() {
	serviceCmd.AddCommand(serviceLogsCmd)

	serviceLogsCmd.Flags().BoolVarP(&flagServiceLogsFollow, "follow", "f", false, "Follow logs")
	serviceLogsCmd.Flags().StringVar(&flagServiceLogsFilter, "filter", "", "Filter pattern to apply")
	serviceLogsCmd.Flags().StringVar(&flagServiceLogsStartTime, "start", "", "Earliest time to return logs (e.g. -1h)")
	serviceLogsCmd.Flags().StringVar(&flagServiceLogsEndTime, "end", "", "Latest time to return logs (e.g. 3y)")
	serviceLogsCmd.Flags().StringSliceVarP(&flagServiceLogsTasks, "task", "t", []string{}, "Show logs from specific task(s)")
}
