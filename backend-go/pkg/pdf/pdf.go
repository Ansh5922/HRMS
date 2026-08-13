package pdf

// PDF generation utility (payslips, offer letters, etc.)
// Uses a third-party library such as gofpdf or chromedp (HTML → PDF)

type Generator struct{}

func (g *Generator) GeneratePayslip(data interface{}) ([]byte, error) {
// TODO: implement HTML template → PDF via headless chrome or gofpdf
return nil, nil
}
