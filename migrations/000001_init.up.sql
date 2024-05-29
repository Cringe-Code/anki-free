create table if not exists users (
    id serial primary key not null,
    name text,
    login text,
    hash_password text
);

create table if not exists packs (
    id serial primary key not null,
    name text,
    rank bigint
);

create table if not exists words (
    id serial primary key not null,
    rus text,
    eng text,
    lvl bigint
);

create table if not exists chances (
    user_id bigint,
    word_id bigint,
    pack_id bigint,
    chance bigint
);

create table if not exists user_pack (
    user_id bigint not null,
    pack_id bigint not null,
    foreign key (user_id) references users(id),
    foreign key (pack_id) references packs(id)
);

create table if not exists pack_word (
    pack_id bigint not null,
    word_id bigint not null,
    foreign key (pack_id) references packs(id),
    foreign key (word_id) references words(id)
);
