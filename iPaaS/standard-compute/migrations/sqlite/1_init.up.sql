create table if not exists job (
    id integer not null default 0 primary key,
    idempotent_id text not null default '',

    state text not null default '',
    sub_state text not null default '',
    file_sync_state text not null default '',

    state_reason text not null default '',
    download_current_size integer default 0,
    download_total_size integer default 0,
    upload_current_size integer default 0,
    upload_total_size integer default 0,

    queue text not null default '',
    priority integer not null default 0,
    request_cores integer not null default 0,
    request_memory integer not null default 0,

    app_mode text not null default '',
    app_path text not null default '',
    singularity_image text not null default '',

    inputs text not null default '',
    output text not null default '',
    env_vars text not null default '',
    command text not null default '',

    workspace text not null default '',
    script text not null default '',
    stdout text not null default '',
    stderr text not null default '',

    is_override integer not null default 0,
    work_dir text not null default '',

    origin_job_id text not null default '',
    alloc_cores integer not null default 0,
    alloc_memory integer not null default 0,
    origin_state text not null default '',
    exit_code text not null default '',
    execution_duration integer not null default 0,
    pending_time text default null,
    running_time text default null,
    completing_time text default null,
    completed_time text default null,

    control_bit_terminate integer not null default 0,
    timeout integer default 0,
    is_timeout integer not null default 0,

    `webhook` text not null default '',

    custom_state_rule text not null default '',

    create_time text not null default '',
    update_time text not null default ''
);
