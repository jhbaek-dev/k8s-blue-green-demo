# Blue-Green Deployment Demo

## 프로젝트 개요

이 프로젝트는 Kubernetes 환경에서 ArgoCD와 Argo Rollouts을 활용한 Blue-Green 배포 전략을 구현하는 데모 시스템입니다. GitOps 방식을 통해 마이크로서비스 애플리케이션의 무중단 배포와 카나리 배포를 안전하게 수행할 수 있는 아키텍처를 제공합니다.

## 아키텍처 구성

### 마이크로서비스 구조

프로젝트는 Go 언어로 작성된 3개의 마이크로서비스로 구성되어 있습니다:

#### Service-A (Frontend Gateway)
- **역할**: 메인 엔트리포인트 역할을 하는 프론트엔드 게이트웨이
- **포트**: 3000 (HTTP)
- **기능**: 
  - Service-B의 HTTP API를 호출하여 버전 정보 수집
  - Service-C의 TCP 소켓에 연결하여 버전 정보 수집
  - 모든 서비스의 버전 정보를 HTML 형태로 통합 제공
  - 환경 변수를 통한 서비스 간 연결 설정
- **특징**: 
  - 다른 서비스와의 연결 상태와 버전 정보를 실시간으로 수집
  - 상세한 디버깅 로그를 통한 서비스 간 통신 추적 가능

#### Service-B (REST API Service)
- **역할**: HTTP REST API를 제공하는 백엔드 서비스
- **포트**: 3000 (HTTP)
- **기능**: 
  - `/version` 엔드포인트를 통한 버전 정보 제공
- **특징**: 
  - 단순하고 경량화된 구조
  - RESTful API 방식의 통신

#### Service-C (TCP Socket Service)
- **역할**: TCP 소켓 통신을 제공하는 백엔드 서비스
- **포트**: 3000 (TCP)
- **기능**: 
  - TCP 연결을 수락하고 버전 정보를 응답
  - 동시 다중 클라이언트 지원 (고루틴 활용)
- **특징**: 
  - Raw TCP 소켓 통신 방식
  - 비동기 처리를 통한 다중 연결 지원

### GitOps 아키텍처

#### App of Apps 패턴
- **Root Application**: `demo-gitops/app-of-apps/root-app.yaml`
  - ArgoCD의 App of Apps 패턴을 구현
  - 하위 애플리케이션들의 배포를 통합 관리
  - Git 리포지토리의 변경사항을 자동으로 감지하여 배포

#### 배포 전략 관리
- **Production Environment**: `prod_version.yaml`
- **QA Environment**: `qa_version.yaml`
- 각 환경별로 독립적인 버전 관리
- 컨테이너 이미지 태그와 루트 버전을 환경별로 분리

### Blue-Green 배포 시스템

#### Argo Rollouts 활용
각 서비스는 Argo Rollouts의 Rollout 리소스를 사용하여 Blue-Green 배포를 구현합니다:

- **자동 프로모션 비활성화**: `autoPromotionEnabled: false`
  - 수동 승인을 통한 안전한 배포 프로세스
  - 배포 전 검증 및 테스트 단계 제공

- **서비스 분리**:
  - **Active Service**: 현재 프로덕션 트래픽을 처리하는 서비스
  - **Preview Service**: 새 버전의 애플리케이션이 배포되는 서비스

#### Istio를 활용한 트래픽 라우팅

Service-B와 Service-C는 Istio VirtualService를 통한 고급 트래픽 라우팅을 지원합니다:

```yaml
# 예시: Service-B VirtualService
- match:
    - sourceLabels:
        root-version: "v2.2.0"
  route:
    - destination:
        host: service-b-preview
```

- **라벨 기반 라우팅**: `root-version` 라벨을 기반으로 트래픽 분기
- **Canary 배포 지원**: 특정 버전의 클라이언트만 preview 환경으로 라우팅
- **점진적 배포**: 트래픽 비율을 조절하여 안전한 배포 진행

## 핵심 기술 스택

### 컨테이너화
- **Docker**: 멀티 스테이지 빌드를 통한 경량화된 컨테이너 이미지
- **Alpine Linux**: 최소한의 베이스 이미지로 보안과 성능 최적화

### 오케스트레이션
- **Kubernetes**: 컨테이너 오케스트레이션 플랫폼
- **ArgoCD**: GitOps 기반 지속적 배포
- **Argo Rollouts**: 고급 배포 전략 (Blue-Green, Canary)

### 서비스 메시
- **Istio**: 서비스 간 통신 관리 및 트래픽 라우팅
- **VirtualService**: 라벨 기반 지능형 트래픽 분기

### 개발 환경
- **Docker Compose**: 로컬 개발 환경 구성
- **Go**: 마이크로서비스 개발 언어

## 배포 프로세스

### 1. 코드 변경 및 이미지 빌드
- 개발자가 코드를 변경하고 Git에 푸시
- CI/CD 파이프라인이 새로운 컨테이너 이미지 빌드
- GitHub Container Registry에 이미지 저장

### 2. GitOps 매니페스트 업데이트
- `prod_version.yaml` 또는 `qa_version.yaml` 파일 업데이트
- 새로운 이미지 태그와 버전 정보 반영

### 3. ArgoCD 자동 배포
- ArgoCD가 Git 리포지토리 변경사항 감지
- App of Apps 패턴을 통해 하위 애플리케이션들 동기화
- Argo Rollouts이 Blue-Green 배포 전략 실행

### 4. 트래픽 라우팅
- Istio VirtualService가 라벨 기반으로 트래픽 분기
- Preview 환경에서 새 버전 검증
- 수동 승인 후 Active 서비스로 트래픽 전환

## 모니터링 및 관찰성

### 버전 추적
- 각 서비스는 자체 버전 정보를 제공
- Service-A를 통해 전체 시스템의 버전 상태 확인 가능
- `root-version` 라벨을 통한 전체 시스템 버전 관리

### 디버깅 지원
- 상세한 로그를 통한 서비스 간 통신 추적
- 연결 상태 및 응답 시간 모니터링
- 환경 변수 설정 상태 확인

## 장점 및 특징

### 무중단 배포
- Blue-Green 배포를 통한 완전한 무중단 서비스
- 배포 실패 시 즉시 롤백 가능
- 사용자 영향 최소화

### 안전한 배포
- 수동 승인 프로세스를 통한 검증 단계
- Preview 환경에서의 사전 테스트
- 라벨 기반 점진적 트래픽 전환

### 확장성
- 마이크로서비스 아키텍처로 독립적인 확장
- Kubernetes의 수평적 확장 지원
- 서비스별 독립적인 배포 사이클

### 관리 효율성
- GitOps 방식의 선언적 인프라 관리
- 버전 관리 시스템을 통한 모든 변경사항 추적
- App of Apps 패턴을 통한 중앙집중식 관리

이 프로젝트는 현대적인 클라우드 네이티브 애플리케이션의 배포 및 운영 모범 사례를 구현한 종합적인 데모 시스템입니다.
