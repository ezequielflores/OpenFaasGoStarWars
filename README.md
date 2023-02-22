# OpenFaas POC

### Ejemplo de implementacion de una funcion Serverless usando el Framework [OpenFaas](https://www.openfaas.com) y [Minikube](https://minikube.sigs.k8s.io/docs/)

Pasos para ejecutar localmente en linux:
1. Instalar [Minikube](https://minikube.sigs.k8s.io/docs/start/)
2. Instalar [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
3. Habilitar el [Docker Registry](https://minikube.sigs.k8s.io/docs/handbook/registry/) local 
4. Instalar herramientas([arkade y faas-cli](https://docs.openfaas.com/cli/install/)) para operar sobre las funciones
5. Instalar [OpenFaas](https://docs.openfaas.com/deployment/kubernetes/#1-deploy-the-chart-with-arkade-fastest-option) usando arkade. _**Es importante prestar atencion a la salida por consola de este comando y seguir las instrucciones.**_