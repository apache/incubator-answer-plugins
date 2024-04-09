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
