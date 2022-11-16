CREATE INDEX IF NOT EXISTS requests_owner_company_idx ON requests(owner_company);
CREATE INDEX IF NOT EXISTS requests_deffect_category_name_idx ON requests(deffect_category_name);
CREATE INDEX IF NOT EXISTS requests_work_type_done_idx ON requests(work_type_done);
CREATE INDEX IF NOT EXISTS requests_serving_company_idx ON requests(serving_company);
CREATE INDEX IF NOT EXISTS request_disp_num_idx ON requests(dispetchers_number);
