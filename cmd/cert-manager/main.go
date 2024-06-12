package main

import (
	"context"
	"os"

	"github.com/glasskube/glasskube/internal/certificates"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	arv1ac "k8s.io/client-go/applyconfigurations/admissionregistration/v1"
	corev1ac "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	serviceName       = "glasskube-webhook-service"
	secretName        = "glasskube-webhook-tls"
	webhookConfigName = "glasskube-validating-webhook-configuration"
	namespace         = "glasskube-system"
	certDir           = ""
	webhookNames      = []string{"vpackage.kb.io", "vclusterpackage.kb.io"}

	log logr.Logger

	fieldManager = "package-operator-cert-manager"
	applyOptions = metav1.ApplyOptions{FieldManager: fieldManager, Force: true}

	cmd = &cobra.Command{
		Use: "cert-manager",
		Run: func(cmd *cobra.Command, args []string) { run(cmd.Context()) },
	}
)

func init() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	log = ctrl.Log.WithName("cert-manager")

	cmd.Flags().StringVar(&certDir, "cert-dir", certDir,
		"directory for certificates (optional)")
	cmd.Flags().StringVar(&serviceName, "service-name", serviceName,
		"name of the webhook service")
	cmd.Flags().StringVar(&secretName, "secret-name", secretName,
		"name of the webhook TLS secret")
	cmd.Flags().StringVar(&webhookConfigName, "webhook-config-name", webhookConfigName,
		"name of the ValidatingWebhookConfiguration to patch")
	cmd.Flags().StringArrayVar(&webhookNames, "webhook-name", webhookNames,
		"name of the webhook to patch")
	cmd.Flags().StringVar(&namespace, "namespace", namespace,
		"namespace of the webhook service and TLS secret")
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Error(err, "command execution failed")
		os.Exit(1)
	}
}

func run(ctx context.Context) {
	client, err := kubernetes.NewForConfig(ctrl.GetConfigOrDie())
	if err != nil {
		log.Error(err, "could not initialize kubernetes client")
		os.Exit(1)
	}

	certificates, err := certificates.Generate(serviceName, namespace, certificates.DefaultValidity)
	if err != nil {
		log.Error(err, "could not generate certificates")
		os.Exit(1)
	}

	webhookEnc, err := certificates.Webhook.Encoded()
	if err != nil {
		log.Error(err, "could not encode certificates")
		os.Exit(1)
	}

	if len(certDir) > 0 {
		if err := webhookEnc.SaveTo(certDir); err != nil {
			log.Error(err, "could not save certificates")
			os.Exit(1)
		}
		log.Info("cerificates saved", "dir", certDir)
	} else {
		secret := corev1ac.Secret(secretName, namespace).
			WithData(webhookEnc.AsMap())
		if _, err := client.CoreV1().Secrets(namespace).Apply(ctx, secret, applyOptions); err != nil {
			log.Error(err, "could not encode certificates", "name", secretName)
			os.Exit(1)
		}
		log.Info("Secret applied", "name", secretName)
	}

	caEnc, err := certificates.Ca.Encoded()
	if err != nil {
		log.Error(err, "could not encode certificates")
		os.Exit(1)
	}

	webhookConfig := arv1ac.ValidatingWebhookConfiguration(webhookConfigName)
	for _, name := range webhookNames {
		webhookConfig.WithWebhooks(
			arv1ac.ValidatingWebhook().
				WithName(name).
				WithClientConfig(arv1ac.WebhookClientConfig().
					WithCABundle(caEnc.Cert...),
				),
		)
	}

	if _, err := client.AdmissionregistrationV1().ValidatingWebhookConfigurations().
		Apply(ctx, webhookConfig, applyOptions); err != nil {
		log.Error(err, "could not apply ValidatingWebhookConfiguration", "name", webhookConfigName)
		os.Exit(1)
	}

	log.Info("ValidatingWebhookConfiguration applied", "name", webhookConfigName)
}
