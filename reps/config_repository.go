package reps

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
)

type ConfigRepository struct {
	Connection *redis.Client
}

func (r *ConfigRepository) GetConfig() (*models.Config, error) {
	var config = &models.Config{}
	conn := r.Connection
	key := "config"
	data, err := conn.Get(context.Background(), key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return config, nil
		} else {
			log.Println("Error getting sources from redis: ", err)
			return nil, err
		}
	}

	err = json.Unmarshal([]byte(data), config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (r *ConfigRepository) SaveConfig(config *models.Config) error {
	conn := r.Connection
	key := "config"
	data, err := json.Marshal(config)
	if err != nil {
		log.Println("Error serializing json: ", err)
		return err
	}

	err = conn.Set(context.Background(), key, data, 0).Err()
	if err != nil {
		log.Println("Error saving ml config to redis: ", err)
		return err
	}

	return nil
}

func (r *ConfigRepository) RestoreConfig() (*models.Config, error) {
	objJson := `{"ai":{"face_recog_mtcnn_threshold":0.86,"face_recog_prob_threshold":0.98,"overlay":true,"plate_recog_instance_count":2,"video_clip_duration":10},"archive":{"action_type":0,"limit_percent":95,"move_location":""},"db":{"connection_string":"mongodb://localhost:27017","type":1},"deep_stack":{"api_key":"","docker_type":1,"fr_enabled":true,"fr_threshold":0.7,"od_enabled":true,"od_threshold":0.45,"performance_mode":1,"server_port":1009,"server_url":"http://127.0.0.1"},"device":{"device_name":"dev","device_type":0},"ffmpeg":{"max_operation_retry_count":10000000,"record_concat_limit":1,"record_video_file_indexer_interval":60,"rtmp_server_init_interval":3.0,"start_task_wait_for_interval":1.0,"use_double_quotes_for_path":false,"watch_dog_failed_wait_interval":3.0,"watch_dog_interval":23},"general":{"heartbeat_interval":30,"root_folder_path":"/mnt/sda1"},"jetson":{"model_name":"ssd-mobilenet-v2"},"jobs":{"mac_ip_matching_enabled":false,"mac_ip_matching_interval":120},"snapshot":{"process_count":4},"source_reader":{"buffer_size":2,"max_retry":150,"max_retry_in":6,"resize_img":false},"tensorflow":{"cache_folder":"/mnt/sdc1/test_projects/tf_cache","model_name":"efficientdet/lite4/detection"},"torch":{"model_name":"ultralytics/yolov5","model_name_specific":"yolov5x6"},"ui":{"booster_interval":0.3,"gs_height":2,"gs_width":4,"seek_to_live_edge_internal":30}}`

	conn := r.Connection
	key := "config"
	err := conn.Set(context.Background(), key, objJson, 0).Err()
	if err != nil {
		log.Println("Error saving default ml config to redis: ", err)
		return nil, err
	}

	var config = &models.Config{}
	err = json.Unmarshal([]byte(objJson), config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
