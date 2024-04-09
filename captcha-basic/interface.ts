/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

export interface FormValue<T = any> {
  value: T;
  isInvalid: boolean;
  errorMsg: string;
  [prop: string]: any;
}

export interface FormDataType {
  [prop: string]: FormValue;
}

export interface FieldError {
  error_field: string;
  error_msg: string;
}

export interface ImgCodeReq {
  captcha_id?: string;
  captcha_code?: string;
}

export interface ImgCodeRes {
  captcha_id: string;
  captcha_img: string;
  verify: boolean;
}

export type CaptchaKey =
  | 'email'
  | 'password'
  | 'edit_userinfo'
  | 'question'
  | 'answer'
  | 'comment'
  | 'edit'
  | 'invitation_answer'
  | 'search'
  | 'report'
  | 'delete'
  | 'vote';
