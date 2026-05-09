package services

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
	"os"

	"chefskiss-backend/models"
)

func SendEmailReceipt(order models.Order) {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_EMAIL")
	senderPass := os.Getenv("SMTP_PASSWORD")

	if senderEmail == "" || senderPass == "" {
		log.Println("Peringatan: Kredensial SMTP belum disetting di .env")
		return
	}

	// 1. Definisikan Template HTML
	const emailHTML = `
	<!DOCTYPE html>
	<html>
	<body style="font-family: 'Inter', sans-serif; background-color: #F9F6F0; margin: 0; padding: 20px;">
		<div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 20px; overflow: hidden; box-shadow: 0 4px 10px rgba(0,0,0,0.05); border: 1px solid #E6C287;">
			<div style="background-color: #8A151B; padding: 30px; text-align: center;">
				<h1 style="color: #ffffff; margin: 0; font-family: 'Playfair Display', serif; font-size: 28px;">Chef's Kiss</h1>
				<p style="color: #E6C287; margin: 5px 0 0 0; font-size: 14px; letter-spacing: 2px;">PREMIUM PASTA & DESSERT</p>
			</div>

			<div style="padding: 40px 30px;">
				<h2 style="color: #2D2D2D; margin-top: 0;">Halo, {{.Customer}}!</h2>
				<p style="color: #555; line-height: 1.6;">Pesananmu telah kami terima dan sedang masuk dalam antrean pre-order. Berikut adalah rincian pesanan kamu:</p>
				
				<div style="background-color: #F9F6F0; border-radius: 12px; padding: 20px; margin: 25px 0;">
					<p style="margin: 0; color: #8A151B; font-weight: bold; font-size: 14px;">TANGGAL PENGAMBILAN:</p>
					<p style="margin: 5px 0 0 0; font-size: 18px; color: #2D2D2D;">{{.PickupDate}}</p>
				</div>

				<table style="width: 100%; border-collapse: collapse;">
					<thead>
						<tr style="border-bottom: 2px solid #F9F6F0;">
							<th style="text-align: left; padding: 10px 0; color: #888; font-size: 12px;">ITEM</th>
							<th style="text-align: center; padding: 10px 0; color: #888; font-size: 12px;">QTY</th>
							<th style="text-align: right; padding: 10px 0; color: #888; font-size: 12px;">HARGA</th>
						</tr>
					</thead>
					<tbody>
						{{range .Items}}
						<tr style="border-bottom: 1px solid #F9F6F0;">
							<td style="padding: 15px 0; color: #2D2D2D; font-weight: 500;">{{.MenuName}}</td>
							<td style="padding: 15px 0; text-align: center; color: #2D2D2D;">{{.Quantity}}x</td>
							<td style="padding: 15px 0; text-align: right; color: #2D2D2D;">Rp {{printf "%.0f" .Price}}</td>
						</tr>
						{{end}}
					</tbody>
				</table>

				<div style="margin-top: 20px; text-align: right;">
					<p style="margin: 0; color: #888; font-size: 14px;">Total Pembayaran:</p>
					<p style="margin: 5px 0 0 0; color: #8A151B; font-size: 24px; font-weight: bold;">Rp {{printf "%.0f" .TotalPrice}}</p>
				</div>
			</div>

			<div style="background-color: #2D2D2D; padding: 30px; text-align: center;">
				<p style="color: #F9F6F0; margin: 0; font-size: 12px; opacity: 0.7;">&copy; 2026 Chef's Kiss</p>
				<p style="color: #E6C287; margin: 10px 0 0 0; font-size: 12px;">Follow us on Instagram <a href="https://instagram.com/chefskisstory" style="color: #E6C287; text-decoration: none; font-weight: bold;">@chefskisstory</a></p>
			</div>
		</div>
	</body>
	</html>
	`

	// 2. Parsing Template dan Inject Data
	t, _ := template.New("receipt").Parse(emailHTML)
	var body bytes.Buffer

	data := struct {
		Customer   string
		PickupDate string
		TotalPrice float64
		Items      []models.OrderItem
	}{
		Customer:   order.Customer,
		PickupDate: order.PickupDate.Format("Monday, 02 January 2006"),
		TotalPrice: order.TotalPrice,
		Items:      order.Items,
	}

	t.Execute(&body, data)

	// 3. Set Header Email (Penting: Content-Type harus text/html)
	subject := "Subject: Konfirmasi Pesanan Chef's Kiss 🍝\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	recipient := order.CustomerEmail
	message := append([]byte(subject+mime), body.Bytes()...)

	// 4. Kirim Email
	auth := smtp.PlainAuth("", senderEmail, senderPass, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{recipient}, message)

	if err != nil {
		log.Println("Gagal mengirim email HTML:", err)
		return
	}

	log.Println("[EMAIL OUTBOX] Sukses mengirim struk HTML ke:", recipient)
}
