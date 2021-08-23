-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------
-- -----------------------------------------------------
-- Schema Locklly
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema Locklly
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `Locklly` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci ;
USE `Locklly` ;

-- -----------------------------------------------------
-- Table `Locklly`.`Address`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `Locklly`.`Address` (
  `ID` VARCHAR(45) NOT NULL,
  `Name` VARCHAR(45) NULL DEFAULT NULL,
  `Address_Line1` VARCHAR(45) NULL DEFAULT NULL,
  `Address_Line2` VARCHAR(45) NULL DEFAULT NULL,
  `City` VARCHAR(45) NULL DEFAULT NULL,
  `Postal_Code` VARCHAR(45) NULL DEFAULT NULL,
  `State` VARCHAR(45) NULL DEFAULT NULL,
  `Phone` VARCHAR(45) NULL DEFAULT NULL,
  PRIMARY KEY (`ID`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `Locklly`.`Brand`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `Locklly`.`Brand` (
  `ID` INT NOT NULL AUTO_INCREMENT,
  `Name` VARCHAR(45) NOT NULL,
  `Description` VARCHAR(45) NULL DEFAULT NULL,
  PRIMARY KEY (`ID`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `Locklly`.`Consumer`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `Locklly`.`Consumer` (
  `ID` INT NOT NULL AUTO_INCREMENT,
  `Is_Active` INT NULL DEFAULT NULL,
  `Address` INT NULL DEFAULT NULL,
  PRIMARY KEY (`ID`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `Locklly`.`BrandMember`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `Locklly`.`BrandMember` (
  `ID` INT NOT NULL AUTO_INCREMENT,
  `Brand_Id` INT NOT NULL,
  `Consumer_Id` INT NOT NULL,
  `Association_Id` VARCHAR(45) NULL DEFAULT NULL,
  `Association_Method` INT NULL DEFAULT NULL,
  PRIMARY KEY (`ID`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `Locklly`.`Category`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `Locklly`.`Category` (
  `ID` VARCHAR(30) NOT NULL,
  `Name` VARCHAR(45) NOT NULL,
  `Description` VARCHAR(45) NULL DEFAULT NULL,
  `ParentId` VARCHAR(45) NULL DEFAULT NULL,
  `PhotoUrl` VARCHAR(500) NULL DEFAULT NULL,
  PRIMARY KEY (`ID`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `Locklly`.`Documents`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `Locklly`.`Documents` (
  `Username` VARCHAR(45) NOT NULL,
  `PANNumber` VARCHAR(45) NULL DEFAULT '""',
  `PANIdURL` VARCHAR(500) NULL DEFAULT '""',
  `BankAccountNumber` VARCHAR(45) NULL DEFAULT '""',
  `BankIFSC` VARCHAR(45) NULL DEFAULT '""',
  `BankIdURL` VARCHAR(500) NULL DEFAULT '""',
  `AadharIdURL` VARCHAR(500) NULL DEFAULT '""',
  `AadharNumber` VARCHAR(45) NULL DEFAULT '""',
  `GSTINNumber` VARCHAR(45) NULL DEFAULT '""',
  `GSTIdURL` VARCHAR(500) NULL DEFAULT '""',
  `PANIdTime` DATETIME NULL DEFAULT NULL,
  `BankIdTime` DATETIME NULL DEFAULT NULL,
  `AadharIdTime` DATETIME NULL DEFAULT NULL,
  `GSTIdTime` DATETIME NULL DEFAULT NULL,
  `PANIdVerified` TINYINT NULL DEFAULT '0',
  `GSTIdVerified` TINYINT NULL DEFAULT '0',
  `BankIdVerified` TINYINT NULL DEFAULT '0',
  `AadharIdVerified` TINYINT NULL DEFAULT '0',
  PRIMARY KEY (`Username`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `Locklly`.`Merchant`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `Locklly`.`Merchant` (
  `Username` VARCHAR(45) NOT NULL,
  `Bio` VARCHAR(50) NULL DEFAULT '""',
  `Active` TINYINT NULL DEFAULT '0',
  `GSTNumber` VARCHAR(45) NULL DEFAULT '""',
  `AadharNumber` VARCHAR(45) NULL DEFAULT '""',
  `PANNumber` VARCHAR(45) NULL DEFAULT '""',
  `PhotoUrl` VARCHAR(500) NULL DEFAULT '""',
  `Name` VARCHAR(45) NULL DEFAULT '""',
  `Verified` TINYINT NULL DEFAULT '0',
  `CreatedAt` DATETIME NULL DEFAULT NULL,
  `FollowersCount` INT NULL DEFAULT '0',
  `ProductViewCount` INT NULL DEFAULT '0',
  `Address` VARCHAR(200) NULL DEFAULT '""',
  `BankAccountNumber` VARCHAR(45) NULL DEFAULT '""',
  `BankIfsc` VARCHAR(45) NULL DEFAULT '""',
  `UPIId` VARCHAR(45) NULL DEFAULT '""',
  `PanCardUrl` VARCHAR(500) NULL DEFAULT '""',
  `AadharCardUrl` VARCHAR(500) NULL DEFAULT '""',
  `Phone` VARCHAR(15) NULL DEFAULT '""',
  `Email` VARCHAR(45) NOT NULL,
  `UpdatedAt` DATETIME NULL DEFAULT NULL,
  PRIMARY KEY (`Username`),
  UNIQUE INDEX `Email_UNIQUE` (`Email` ASC) ,
  INDEX `Address_idx` (`Bio` ASC) )
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `Locklly`.`Offer`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `Locklly`.`Offer` (
  `Id` INT NOT NULL,
  `Code` VARCHAR(45) NULL DEFAULT NULL,
  `Name` VARCHAR(45) NULL DEFAULT NULL,
  `Description` VARCHAR(45) NULL DEFAULT NULL,
  PRIMARY KEY (`Id`),
  UNIQUE INDEX `Code_UNIQUE` (`Code` ASC) )
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `Locklly`.`Product`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `Locklly`.`Product` (
  `ID` VARCHAR(10) NOT NULL,
  `Title` VARCHAR(45) NOT NULL,
  `Description` VARCHAR(45) NULL DEFAULT NULL,
  `CreatedBy` VARCHAR(45) NOT NULL,
  `CreatedAt` DATETIME NULL DEFAULT NULL,
  `PhotoUrls` VARCHAR(5000) NULL DEFAULT NULL,
  `Category` VARCHAR(100) NULL DEFAULT NULL,
  `Price` DECIMAL(7,2) NULL DEFAULT '0.00',
  `Likes` INT NULL DEFAULT '0',
  `Rating` DECIMAL(3,2) NULL DEFAULT '0.00',
  `Sizes` VARCHAR(50) NULL DEFAULT NULL,
  `CountPeopleRated` INT NULL DEFAULT '0',
  `PrimaryPhotoUrl` VARCHAR(500) NULL DEFAULT NULL,
  `InterestingFact` VARCHAR(200) NULL DEFAULT NULL,
  `Views` INT NULL DEFAULT NULL,
  `Tags` VARCHAR(100) NULL DEFAULT NULL,
  `Events` VARCHAR(100) NULL DEFAULT NULL,
  PRIMARY KEY (`ID`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `Locklly`.`Tags`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `Locklly`.`Tags` (
  `Name` VARCHAR(45) NOT NULL,
  `Description` VARCHAR(45) NULL DEFAULT NULL,
  PRIMARY KEY (`Name`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `Locklly`.`TestTable`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `Locklly`.`TestTable` (
  `one` VARCHAR(10) NULL DEFAULT NULL,
  `two` VARCHAR(10) NULL DEFAULT NULL)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `Locklly`.`Users`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `Locklly`.`Users` (
  `UserName` VARCHAR(30) NOT NULL,
  `PhoneNumber` VARCHAR(45) NULL DEFAULT NULL,
  `EmailId` VARCHAR(45) NULL DEFAULT NULL,
  `SignupMethod` INT NOT NULL,
  PRIMARY KEY (`UserName`),
  UNIQUE INDEX `UserName_UNIQUE` (`UserName` ASC) ,
  UNIQUE INDEX `Credential_UNIQUE` (`PhoneNumber` ASC, `EmailId` ASC) )
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


-- -----------------------------------------------------
-- Table `Locklly`.`UserInfo`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `Locklly`.`UserInfo` (
  `UserName` VARCHAR(30) NOT NULL,
  `Name` VARCHAR(45) NULL DEFAULT NULL,
  `AboutLine` VARCHAR(45) NULL DEFAULT NULL,
  `PhotoUrl` VARCHAR(500) NULL DEFAULT NULL,
  `City` VARCHAR(45) NULL DEFAULT NULL,
  `Region` VARCHAR(45) NULL DEFAULT NULL,
  `Gender` VARCHAR(10) NULL DEFAULT NULL,
  `State` VARCHAR(45) NULL DEFAULT NULL,
  PRIMARY KEY (`UserName`),
  UNIQUE INDEX `PhotoUrl_UNIQUE` (`PhotoUrl` ASC) )
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8
COLLATE = utf8_general_ci;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
