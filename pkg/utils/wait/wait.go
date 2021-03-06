/*
 * Tencent is pleased to support the open source community by making TKEStack available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package wait

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

// RunUntil run fc period until ctx is done
func RunUntil(ctx context.Context, log logrus.FieldLogger, interval time.Duration, fc func() error) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if err := fc(); err != nil {
			log.Errorf(err.Error())
		}
		time.Sleep(interval)
	}
}

// RunEvent run fc when event chan receive data
func RunEvent(ctx context.Context, event chan struct{}, fc func()) {
	for {
		select {
		case <-ctx.Done():
		case <-event:
			fc()
		}
	}
}
