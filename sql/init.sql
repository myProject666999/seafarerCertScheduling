CREATE DATABASE IF NOT EXISTS seafarer_cert_scheduling
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;

USE seafarer_cert_scheduling;

-- ============================================================
-- 船员信息表
-- ============================================================
CREATE TABLE seafarer (
    id            BIGINT       NOT NULL AUTO_INCREMENT,
    name          VARCHAR(50)  NOT NULL,
    gender        TINYINT      NOT NULL DEFAULT 0 COMMENT '0-未知 1-男 2-女',
    birthday      DATE                  DEFAULT NULL,
    id_number     VARCHAR(18)           DEFAULT NULL,
    phone         VARCHAR(20)           DEFAULT NULL,
    email         VARCHAR(100)          DEFAULT NULL,
    rank          VARCHAR(50)           DEFAULT NULL COMMENT '职务/职级',
    status        TINYINT      NOT NULL DEFAULT 0 COMMENT '0-待派 1-在船 2-休假',
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at    DATETIME              DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_seafarer_status (status),
    INDEX idx_seafarer_id_number (id_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='船员信息表';

-- ============================================================
-- 证书类型表
-- ============================================================
CREATE TABLE certificate_type (
    id               BIGINT       NOT NULL AUTO_INCREMENT,
    name             VARCHAR(100) NOT NULL,
    code             VARCHAR(50)  NOT NULL,
    description      TEXT                  DEFAULT NULL,
    validity_months  INT                   DEFAULT NULL COMMENT '有效月数 NULL=长期有效',
    is_required      TINYINT      NOT NULL DEFAULT 0 COMMENT '0-否 1-是',
    created_at       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_cert_type_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='证书类型表';

-- ============================================================
-- 船员证书表（一张证只归一个船员）
-- ============================================================
CREATE TABLE seafarer_certificate (
    id                   BIGINT       NOT NULL AUTO_INCREMENT,
    seafarer_id          BIGINT       NOT NULL,
    certificate_type_id  BIGINT       NOT NULL,
    cert_number          VARCHAR(100) NOT NULL,
    issue_date           DATE         NOT NULL,
    expire_date          DATE                  DEFAULT NULL COMMENT 'NULL=长期有效',
    cert_image_url       VARCHAR(500)          DEFAULT NULL,
    status               TINYINT      NOT NULL DEFAULT 1 COMMENT '0-已过期 1-有效 2-即将过期',
    created_at           DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at           DATETIME              DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX idx_seafarer_cert_seafarer (seafarer_id),
    INDEX idx_seafarer_cert_expire (expire_date),
    INDEX idx_seafarer_cert_status (status),
    CONSTRAINT fk_seafarer_cert_seafarer FOREIGN KEY (seafarer_id) REFERENCES seafarer(id),
    CONSTRAINT fk_seafarer_cert_type     FOREIGN KEY (certificate_type_id) REFERENCES certificate_type(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='船员证书表';

-- ============================================================
-- 船舶信息表
-- ============================================================
CREATE TABLE ship (
    id            BIGINT        NOT NULL AUTO_INCREMENT,
    name          VARCHAR(100)  NOT NULL,
    imo_number    VARCHAR(7)             DEFAULT NULL,
    mmsi          VARCHAR(9)             DEFAULT NULL,
    ship_type     VARCHAR(50)            DEFAULT NULL,
    gross_tonnage DECIMAL(10,2)          DEFAULT NULL,
    flag_state    VARCHAR(50)            DEFAULT NULL,
    status        TINYINT       NOT NULL DEFAULT 1 COMMENT '0-报废 1-运营 2-维修',
    created_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at    DATETIME               DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_ship_imo (imo_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='船舶信息表';

-- ============================================================
-- 船舶岗位编制表
-- ============================================================
CREATE TABLE ship_position (
    id            BIGINT      NOT NULL AUTO_INCREMENT,
    ship_id       BIGINT      NOT NULL,
    position_name VARCHAR(50) NOT NULL,
    department    VARCHAR(50)          DEFAULT NULL COMMENT '甲板部/轮机部/客运部',
    required_count INT        NOT NULL DEFAULT 1,
    sort_order    INT         NOT NULL DEFAULT 0,
    created_at    DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_ship_position_ship (ship_id),
    CONSTRAINT fk_ship_position_ship FOREIGN KEY (ship_id) REFERENCES ship(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='船舶岗位编制表';

-- ============================================================
-- 岗位证书要求表
-- ============================================================
CREATE TABLE position_cert_requirement (
    id                  BIGINT  NOT NULL AUTO_INCREMENT,
    ship_position_id    BIGINT  NOT NULL,
    certificate_type_id BIGINT  NOT NULL,
    is_mandatory        TINYINT NOT NULL DEFAULT 1 COMMENT '0-否 1-是',
    created_at          DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_position_cert (ship_position_id, certificate_type_id),
    CONSTRAINT fk_pos_cert_position FOREIGN KEY (ship_position_id) REFERENCES ship_position(id),
    CONSTRAINT fk_pos_cert_type     FOREIGN KEY (certificate_type_id) REFERENCES certificate_type(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='岗位证书要求表';

-- ============================================================
-- 船员在船分配表
-- ============================================================
CREATE TABLE seafarer_assignment (
    id                       BIGINT  NOT NULL AUTO_INCREMENT,
    seafarer_id              BIGINT  NOT NULL,
    ship_id                  BIGINT  NOT NULL,
    ship_position_id         BIGINT  NOT NULL,
    embark_date              DATE    NOT NULL,
    expected_disembark_date  DATE             DEFAULT NULL,
    actual_disembark_date    DATE             DEFAULT NULL,
    status                   TINYINT NOT NULL DEFAULT 1 COMMENT '0-已下船 1-在船',
    created_at               DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_assignment_seafarer (seafarer_id),
    INDEX idx_assignment_ship (ship_id),
    INDEX idx_assignment_position (ship_position_id),
    INDEX idx_assignment_status (status),
    CONSTRAINT fk_assignment_seafarer FOREIGN KEY (seafarer_id) REFERENCES seafarer(id),
    CONSTRAINT fk_assignment_ship     FOREIGN KEY (ship_id) REFERENCES ship(id),
    CONSTRAINT fk_assignment_position FOREIGN KEY (ship_position_id) REFERENCES ship_position(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='船员在船分配表';

-- ============================================================
-- 航次合同表
-- ============================================================
CREATE TABLE voyage_contract (
    id               BIGINT       NOT NULL AUTO_INCREMENT,
    seafarer_id      BIGINT       NOT NULL,
    ship_id          BIGINT       NOT NULL,
    contract_number  VARCHAR(50)  NOT NULL,
    start_date       DATE         NOT NULL,
    end_date         DATE         NOT NULL,
    actual_end_date  DATE                  DEFAULT NULL,
    status           TINYINT      NOT NULL DEFAULT 1 COMMENT '0-已终止 1-执行中 2-已完成',
    remarks          TEXT                  DEFAULT NULL,
    created_at       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_contract_seafarer (seafarer_id),
    INDEX idx_contract_ship (ship_id),
    INDEX idx_contract_number (contract_number),
    CONSTRAINT fk_contract_seafarer FOREIGN KEY (seafarer_id) REFERENCES seafarer(id),
    CONSTRAINT fk_contract_ship     FOREIGN KEY (ship_id) REFERENCES ship(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='航次合同表';

-- ============================================================
-- 上下船记录表
-- ============================================================
CREATE TABLE embark_disembark_record (
    id           BIGINT       NOT NULL AUTO_INCREMENT,
    seafarer_id  BIGINT       NOT NULL,
    ship_id      BIGINT       NOT NULL,
    record_type  TINYINT      NOT NULL COMMENT '1-上船 2-下船',
    record_date  DATE         NOT NULL,
    port         VARCHAR(100)          DEFAULT NULL,
    reason       VARCHAR(200)          DEFAULT NULL,
    operator     VARCHAR(50)           DEFAULT NULL,
    created_at   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_embark_seafarer (seafarer_id),
    INDEX idx_embark_ship (ship_id),
    INDEX idx_embark_date (record_date),
    CONSTRAINT fk_embark_seafarer FOREIGN KEY (seafarer_id) REFERENCES seafarer(id),
    CONSTRAINT fk_embark_ship     FOREIGN KEY (ship_id) REFERENCES ship(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='上下船记录表';

-- ============================================================
-- 休假记录表
-- ============================================================
CREATE TABLE leave_record (
    id           BIGINT      NOT NULL AUTO_INCREMENT,
    seafarer_id  BIGINT      NOT NULL,
    start_date   DATE        NOT NULL,
    end_date     DATE                 DEFAULT NULL,
    leave_days   INT                  DEFAULT NULL,
    status       TINYINT     NOT NULL DEFAULT 1 COMMENT '0-已结束 1-休假中',
    reason       VARCHAR(200)         DEFAULT NULL,
    created_at   DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_leave_seafarer (seafarer_id),
    INDEX idx_leave_status (status),
    CONSTRAINT fk_leave_seafarer FOREIGN KEY (seafarer_id) REFERENCES seafarer(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='休假记录表';

-- ============================================================
-- 健康复检记录表
-- ============================================================
CREATE TABLE health_reexamination (
    id               BIGINT       NOT NULL AUTO_INCREMENT,
    seafarer_id      BIGINT       NOT NULL,
    exam_date        DATE         NOT NULL,
    next_exam_date   DATE                  DEFAULT NULL,
    exam_result      TINYINT      NOT NULL COMMENT '1-合格 2-不合格 3-限制',
    exam_institution VARCHAR(200)          DEFAULT NULL,
    report_url       VARCHAR(500)          DEFAULT NULL,
    restrictions     TEXT                  DEFAULT NULL,
    created_at       DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_health_seafarer (seafarer_id),
    INDEX idx_health_next_exam (next_exam_date),
    CONSTRAINT fk_health_seafarer FOREIGN KEY (seafarer_id) REFERENCES seafarer(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='健康复检记录表';

-- ============================================================
-- 船员调动申请表
-- ============================================================
CREATE TABLE transfer_request (
    id                     BIGINT       NOT NULL AUTO_INCREMENT,
    seafarer_id            BIGINT       NOT NULL,
    from_ship_id           BIGINT       NOT NULL,
    to_ship_id             BIGINT       NOT NULL,
    from_position_id       BIGINT       NOT NULL,
    to_position_id         BIGINT       NOT NULL,
    replacement_seafarer_id BIGINT               DEFAULT NULL COMMENT '原船补位船员',
    reason                 VARCHAR(500) NOT NULL,
    status                 TINYINT      NOT NULL DEFAULT 0 COMMENT '0-待审批 1-已批准 2-已拒绝 3-已取消',
    approver               VARCHAR(50)           DEFAULT NULL,
    approve_remark         VARCHAR(500)          DEFAULT NULL,
    approved_at            DATETIME              DEFAULT NULL,
    from_ship_valid        TINYINT               DEFAULT NULL COMMENT '原船编制校验 0-不通过 1-通过',
    to_ship_valid          TINYINT               DEFAULT NULL COMMENT '目标船编制校验 0-不通过 1-通过',
    created_at             DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at             DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_transfer_seafarer (seafarer_id),
    INDEX idx_transfer_from_ship (from_ship_id),
    INDEX idx_transfer_to_ship (to_ship_id),
    INDEX idx_transfer_status (status),
    CONSTRAINT fk_transfer_seafarer      FOREIGN KEY (seafarer_id) REFERENCES seafarer(id),
    CONSTRAINT fk_transfer_from_ship     FOREIGN KEY (from_ship_id) REFERENCES ship(id),
    CONSTRAINT fk_transfer_to_ship       FOREIGN KEY (to_ship_id) REFERENCES ship(id),
    CONSTRAINT fk_transfer_from_position FOREIGN KEY (from_position_id) REFERENCES ship_position(id),
    CONSTRAINT fk_transfer_to_position   FOREIGN KEY (to_position_id) REFERENCES ship_position(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='船员调动申请表';

-- ============================================================
-- 证书预警记录表
-- ============================================================
CREATE TABLE cert_alert (
    id                      BIGINT  NOT NULL AUTO_INCREMENT,
    seafarer_certificate_id BIGINT  NOT NULL,
    seafarer_id             BIGINT  NOT NULL,
    alert_level             TINYINT NOT NULL COMMENT '1-90天 2-60天 3-30天',
    alert_date              DATE    NOT NULL,
    expire_date             DATE    NOT NULL,
    days_remaining          INT     NOT NULL,
    is_handled              TINYINT NOT NULL DEFAULT 0 COMMENT '0-否 1-是',
    handle_remark           VARCHAR(500)     DEFAULT NULL,
    created_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    INDEX idx_alert_seafarer (seafarer_id),
    INDEX idx_alert_level (alert_level),
    INDEX idx_alert_handled (is_handled),
    INDEX idx_alert_date (alert_date),
    CONSTRAINT fk_alert_cert     FOREIGN KEY (seafarer_certificate_id) REFERENCES seafarer_certificate(id),
    CONSTRAINT fk_alert_seafarer FOREIGN KEY (seafarer_id) REFERENCES seafarer(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='证书预警记录表';

-- ============================================================
-- 初始数据：常见海船证书类型
-- ============================================================
INSERT INTO certificate_type (name, code, validity_months, is_required) VALUES
('GMDSS操作员证书',             'GMDSS',       60, 1),
('精通救生艇筏和救助艇',         'PROF_LB',     60, 1),
('高级消防',                     'ADV_FF',      60, 1),
('精通急救',                     'MED_FA',      60, 1),
('船舶保安员',                   'SSO',         60, 1),
('ISPS保安意识培训',             'ISPS_AW',     60, 0),
('海船船员健康证',               'HEALTH',      24, 1),
('船员服务簿',                   'SRV_BOOK',    NULL, 1),
('适任证书（船长）',             'CO_CAPTAIN',  60, 1),
('适任证书（大副）',             'CO_CHIEF',    60, 1),
('适任证书（二副）',             'CO_2ND',      60, 1),
('适任证书（三副）',             'CO_3RD',      60, 1),
('适任证书（轮机长）',           'CO_CENG',     60, 1),
('适任证书（大管轮）',           'CO_1AE',      60, 1),
('适任证书（二管轮）',           'CO_2AE',      60, 1),
('适任证书（三管轮）',           'CO_3AE',      60, 1),
('适任证书（水手）',             'CO_SAILOR',   60, 0),
('适任证书（机工）',             'CO_MOTORMAN', 60, 0),
('油船货物操作高级培训',         'OIL_ADV',     60, 0),
('化学品船货物操作高级培训',     'CHEM_ADV',    60, 0);
