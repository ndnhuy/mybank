create database if not exists mybankdb;

create table if not exists accounts(
  id char(36) not null,
  balance decimal(10, 2),
  primary key(`id`)
) engine = InnoDB;