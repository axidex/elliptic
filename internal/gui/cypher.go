package gui

import (
	"crypto/rand"
	"errors"
	"fyne.io/fyne/v2/dialog"
	"github.com/axidex/elliptic/internal/cypher"
)

var (
	ErrGeneratingKeys = errors.New("error generating keys")
	ErrImportingKeys  = errors.New("error importing keys")
	ErrEncrypt        = errors.New("error encrypt")
	ErrMessage        = errors.New("empty message")
	ErrDecrypt        = errors.New("error decrypt")
)

func (app *AppGui) generateKeys() {
	app.logger.Infof("Generating keys")

	if app.publicKeyEntry == nil || app.privateKeyEntry == nil {
		app.logger.Errorf("Public or private key is nil")
		dialog.ShowError(ErrGeneratingKeys, app.w)
		return
	}
	//params := ParamName(app.selectParams.Selected).GetParamByName()
	curve := CurveName(app.selectCurve.Selected).GetCurveByName()
	// Генерируем случайные ключи
	keys, err := cypher.GenerateKey(rand.Reader, curve, nil)
	if err != nil {
		app.logger.Errorf("GenerateKey err: %s", err)
		dialog.ShowError(ErrGeneratingKeys, app.w)
		return
	}

	app.setPrivateKeyInfo(keys)

	private, err := cypher.ExportPrivatePEM(keys)
	if err != nil {
		app.logger.Errorf("Encoding err: %s", err)
		dialog.ShowError(ErrGeneratingKeys, app.w)
		return
	}

	public, err := cypher.ExportPublicPEM(&keys.PublicKey)
	if err != nil {
		app.logger.Errorf("Encoding err: %s", err)
		dialog.ShowError(ErrGeneratingKeys, app.w)
		return
	}

	// Обновляем текст в окнах с ключами
	app.privateKeyEntry.SetText(string(private))
	app.publicKeyEntry.SetText(string(public))
}

func (app *AppGui) encryptData() {
	text := app.openText.Text
	//if len(text) == 0 {
	//	dialog.ShowError(ErrMessage, app.w)
	//	return
	//}

	pemKey := []byte(app.publicKeyEntry.Text)

	app.logger.Infof("Got task encryption")
	key, err := cypher.ImportPublicPEM(pemKey)
	if err != nil {
		app.logger.Infof("Not valid key: %v", err)
		dialog.ShowError(ErrImportingKeys, app.w)
		return
	}

	app.setPublicKeyInfo(key)

	encryptedText, err := cypher.Encrypt(rand.Reader, key, []byte(text), nil, nil)
	if err != nil {
		app.logger.Errorf("Encryption error %v", err)
		dialog.ShowError(err, app.w)
		return
	}

	app.closedText.SetText(encryptedText)
}

func (app *AppGui) decryptData() {
	encryptedBytes := app.closedText.Text

	pemKey := []byte(app.privateKeyEntry.Text)

	app.logger.Infof("Got task decryption")

	keys, err := cypher.ImportPrivatePEM(pemKey)
	if err != nil {
		app.logger.Infof("Not valid key: %v", err)
		dialog.ShowError(ErrImportingKeys, app.w)
		return
	}

	app.setPrivateKeyInfo(keys)

	//app.logger.Infof("Decrypting data %s", encryptedBytes)

	decryptText, err := keys.Decrypt(rand.Reader, encryptedBytes, nil, nil)
	if err != nil {
		app.logger.Infof("Decryption error %v", err)
		dialog.ShowError(err, app.w)
		return
	}

	app.openText.SetText(string(decryptText))
}
