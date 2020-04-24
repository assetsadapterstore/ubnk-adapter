/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package openwtester

import (
	"github.com/blocktree/openwallet/v2/openw"
	"testing"
	"time"

	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
)

func testGetAssetsAccountBalance(tm *openw.WalletManager, walletID, accountID string) {
	balance, err := tm.GetAssetsAccountBalance(testApp, walletID, accountID)
	if err != nil {
		log.Error("GetAssetsAccountBalance failed, unexpected error:", err)
		return
	}
	log.Info("balance:", balance)
}

func testGetAssetsAccountTokenBalance(tm *openw.WalletManager, walletID, accountID string, contract openwallet.SmartContract) {
	balance, err := tm.GetAssetsAccountTokenBalance(testApp, walletID, accountID, contract)
	if err != nil {
		log.Error("GetAssetsAccountTokenBalance failed, unexpected error:", err)
		return
	}
	log.Info("token balance:", balance.Balance)
}

func testCreateTransactionStep(tm *openw.WalletManager, walletID, accountID, to, amount, feeRate string, contract *openwallet.SmartContract) (*openwallet.RawTransaction, error) {

	//err := tm.RefreshAssetsAccountBalance(testApp, accountID)
	//if err != nil {
	//	log.Error("RefreshAssetsAccountBalance failed, unexpected error:", err)
	//	return nil, err
	//}

	rawTx, err := tm.CreateTransaction(testApp, walletID, accountID, amount, to, feeRate, "", contract)

	if err != nil {
		log.Error("CreateTransaction failed, unexpected error:", err)
		return nil, err
	}

	return rawTx, nil
}

func testCreateSummaryTransactionStep(
	tm *openw.WalletManager,
	walletID, accountID, summaryAddress, minTransfer, retainedBalance, feeRate string,
	start, limit int,
	contract *openwallet.SmartContract) ([]*openwallet.RawTransaction, error) {

	rawTxArray, err := tm.CreateSummaryTransaction(testApp, walletID, accountID, summaryAddress, minTransfer,
		retainedBalance, feeRate, start, limit, contract)

	if err != nil {
		log.Error("CreateSummaryTransaction failed, unexpected error:", err)
		return nil, err
	}

	return rawTxArray, nil
}

func testSignTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	_, err := tm.SignTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, "12345678", rawTx)
	if err != nil {
		log.Error("SignTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Infof("rawTx: %+v", rawTx)
	return rawTx, nil
}

func testVerifyTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	//log.Info("rawTx.Signatures:", rawTx.Signatures)

	_, err := tm.VerifyTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("VerifyTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Infof("rawTx: %+v", rawTx)
	return rawTx, nil
}

func testSubmitTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	tx, err := tm.SubmitTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("SubmitTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Std.Info("tx: %+v", tx)
	log.Info("wxID:", tx.WxID)
	log.Info("txID:", rawTx.TxID)

	return rawTx, nil
}

func TestTransfer(t *testing.T) {

	addrs := []string{
		"02116c41043ea98d39fa3a651530f8c5c389821b8bb3337fb246690b091245996f",
		"02462e1a68d9f11a75b3aa1967105ded98f12ac1e0b4a65cab3919469bd367a838",
		"02597d17078b4282f225893d7c293ab24e952d5dbf915756c52ec86eda1b5df6d3",
		"02aa0c834992e2e89a49138dbf32e1ef7dc8f0462f976583d9c2a8da0311b92ab6",
		"02adb86d29c6158cd2e114a1d53815f77421c2c00b4dc2a84be0c6e6bf596565ad",
		"02ddbc390f30709725d6d9ac07fca13393dbcfbe1c1f3cadbb6d401106454b624a",
		"03168332bc84bbff92a5e96fd5cb725468fe823982ffac90719b520eae19f0adfa",
		"031b0fd4cf533122efafdea5245485139af2cde74c8a6ec164ae5d1fd1127531c3",
		"031ea240cace57b2148d12b0dddd8dcfe6b5a924eaa6deaccb360508a5e89c210d",
		"033a4c68edfb336928ff5582feb4c65c08011369d1862f5af1a7f5d8cb3e1a313a",
		"03f2b9f83f2ea015453e81561cd1d0a4af0f9f748b3175bb0188c87531edb82b3f",

	}

	tm := testInitWalletManager()
	walletID := "WGE8SRfrTTZCceMx5W3VtBVaUecCGSpg5N"
	accountID := "YnsCLKL54J9qBPMaBHsALvYySQcY3eqod1wQeBRDugD"

	testGetAssetsAccountBalance(tm, walletID, accountID)

	for _, to := range addrs {

		rawTx, err := testCreateTransactionStep(tm, walletID, accountID, to, "1", "", nil)
		if err != nil {
			return
		}

		_, err = testSignTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testVerifyTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testSubmitTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		time.Sleep(5 * time.Second)
	}
}

func TestSummary(t *testing.T) {
	tm := testInitWalletManager()
	walletID := "WGE8SRfrTTZCceMx5W3VtBVaUecCGSpg5N"
	accountID := "2Tk6kPzPuv3x97huu9r8GDXTRexAGYPCjfTinbfhxEmC"
	summaryAddress := "02551cd92569da93f35aa91db48c7bfb1eff7f673d7ce1aa544eab74844eaf9748"

	testGetAssetsAccountBalance(tm, walletID, accountID)

	rawTxArray, err := testCreateSummaryTransactionStep(tm, walletID, accountID,
		summaryAddress, "", "", "",
		0, 100, nil)
	if err != nil {
		log.Errorf("CreateSummaryTransaction failed, unexpected error: %v", err)
		return
	}

	//执行汇总交易
	for _, rawTx := range rawTxArray {
		_, err = testSignTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testVerifyTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testSubmitTransactionStep(tm, rawTx)
		if err != nil {
			return
		}
	}

}
